package incoming_M1

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/lent0s/asyncHTTPhandlerOnRabbitMQ/rabbit_MQ"
	"io"
	"log"
	"net/http"
	"time"
)

func (c *control) worker(w http.ResponseWriter, r *http.Request) {

	switch {
	case r.Method != http.MethodPost:
		notSupportedMethod(w, r.Method)
		return
	case c.interrupt:
		http.Error(w, "service unavailable", http.StatusServiceUnavailable)
		return
	}

	c.clientsON.Add(1)
	defer func() {
		c.clientsON.Done()
	}()

	body, err := getBody(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	id := makeID(body)

	if err = rabbit_MQ.PostM1(
		append([]byte(id), body...),
		r.Header.Get("Content-Type"),
		id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	body, contType := rabbit_MQ.GetFromM2(id)
	w.Header().Add("Content-Type", contType)
	b, err := w.Write(body)
	if err != nil {
		log.Printf("id [%s] response interruption [%d/%d]: %s\n",
			id, b, len(body), err)
		w.WriteHeader(http.StatusPartialContent)
		return
	}
	log.Printf("[%s]<- (Got) %s\n", id, body)
}

func makeID(data []byte) string {

	data = append([]byte(time.Now().String()), data...)
	hash := sha256.Sum256(data)
	return hex.EncodeToString(hash[:])[:rabbit_MQ.LenID]
}

func getBody(body io.ReadCloser) ([]byte, error) {

	data, err := io.ReadAll(body)
	if err != nil {
		return nil, err
	}
	_ = body.Close()

	return data, nil
}

func notSupportedMethod(w http.ResponseWriter, method string) {

	w.Header().Set("Support", "[GET/POST]")
	err := fmt.Sprintf("method \"%s\" not allowed. Need [GET/POST]", method)
	http.Error(w, err, http.StatusMethodNotAllowed)
}
