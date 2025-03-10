package actions

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"grafana-matrix-alerts/logger"
	"grafana-matrix-alerts/mautrixClient"
	"io"
	"maunium.net/go/mautrix/event"
	"maunium.net/go/mautrix/format"
	"maunium.net/go/mautrix/id"
	"strings"
)

func HandleWebhook(c *gin.Context) {

	jsonData, err := io.ReadAll(c.Request.Body)
	if err != nil {
		logger.Log.Error().
			Err(err).
			Msg("Unable to read request body")

		return
	}

	var data AlertsData
	if err := json.Unmarshal(jsonData, &data); err != nil {
		logger.Log.Error().
			Err(err).
			Msg("Unable to unmarshal request body")

		return
	}

	// Range trough alerts
	for _, alertEntry := range data.Alerts {

		var alertAction string
		switch alertEntry.Status {
		case "firing":
			alertAction = "💔 ALERT"
		case "resolved":
			alertAction = "💚 RESOLVED"
		default:
			alertAction = "❓ NO DATA"
		}

		alertMessage := fmt.Sprintf(
			"<b>%s</b> <br> Rule: <a href='%s'>%s</a> | %s<br><br> <pre><code>%s</code></pre>",
			alertAction,
			alertEntry.DashboardURL,
			alertEntry.Labels.Alertname,
			alertEntry.Annotations.Summary,
			strings.ReplaceAll(alertEntry.ValueString, "], [", "],\n["),
		)

		content := format.RenderMarkdown(alertMessage, true, true)
		_, err = mautrixClient.Client.SendMessageEvent(context.TODO(), id.RoomID(c.Param("room")), event.EventMessage, &content)
		if err != nil {
			logger.Log.Error().
				Err(err).
				Msg("Unable to send message event")
		}
	}

}
