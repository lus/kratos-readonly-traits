package webhook

import (
	"encoding/json"
	"github.com/lus/kratos-readonly-traits/internal/schema"
	"github.com/rs/zerolog/log"
	"io"
	"net/http"
)

type requestPayload struct {
	SchemaURL string         `json:"schema_url"`
	OldTraits map[string]any `json:"old_traits"`
	NewTraits map[string]any `json:"new_traits"`
}

type responsePayload struct {
	Messages []responseTopMessage `json:"messages"`
}

type responseTopMessage struct {
	InstancePtr string                  `json:"instance_ptr"`
	Messages    []responseNestedMessage `json:"messages"`
}

type responseNestedMessage struct {
	ID   int    `json:"id"`
	Text string `json:"text"`
	Type string `json:"type"`
}

type controller struct {
	ErrorMessage string
}

func (cnt *controller) Endpoint(writer http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodPost {
		writer.WriteHeader(http.StatusMethodNotAllowed)
		writer.Write([]byte("method not allowed"))
		return
	}

	body, err := io.ReadAll(request.Body)
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		writer.Write([]byte("could not read request body"))
		return
	}

	data := new(requestPayload)
	if err := json.Unmarshal(body, data); err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		writer.Write([]byte("could not unmarshal request body"))
		return
	}

	traits, err := schema.ExtractReadOnlyTraits(data.SchemaURL)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		writer.Write([]byte(err.Error()))
		return
	}

	response := &responsePayload{
		Messages: make([]responseTopMessage, 0, len(traits)),
	}
	for trait, readonly := range traits {
		if !readonly {
			continue
		}
		oldValue, ok := data.OldTraits[trait]
		if !ok {
			continue
		}
		newValue, ok := data.NewTraits[trait]
		if !ok {
			continue
		}
		if oldValue != newValue {
			log.Debug().
				Str("trait", trait).
				Interface("old_value", oldValue).
				Interface("new_value", newValue).
				Msg("Read-only trait changed.")
			response.Messages = append(response.Messages, responseTopMessage{
				InstancePtr: "#/traits/" + trait,
				Messages: []responseNestedMessage{
					{
						ID:   1337,
						Text: cnt.ErrorMessage,
						Type: "conflict",
					},
				},
			})
		}
	}
	if len(response.Messages) > 0 {
		responseData, err := json.Marshal(response)
		if err != nil {
			writer.WriteHeader(http.StatusInternalServerError)
			writer.Write([]byte(err.Error()))
			return
		}
		writer.Header().Set("Content-Type", "application/json")
		writer.WriteHeader(http.StatusConflict)
		writer.Write(responseData)
		return
	}
	writer.WriteHeader(http.StatusOK)
}
