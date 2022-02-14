package api

import (
	"fmt"
	"math/rand"
	"net/http"

	"github.com/AgoraIO-Community/Cloud-Recording-Golang/schemas"
	"github.com/AgoraIO-Community/Cloud-Recording-Golang/utils"
	"github.com/gofiber/fiber/v2"
)

func startRecording(c *fiber.Ctx) error {
	u := new(schemas.CallInfo)

	if err := c.BodyParser(u); err != nil {
		return c.Status(http.StatusUnprocessableEntity).JSON(fiber.Map{
			"msg": "invalid json",
			"err": err.Error(),
		})
	}
	uid := int(rand.Uint32())
	rec := &utils.Recorder{
		CallInfo: *u,
		UID:      uid,
	}

	_, err := rec.Acquire()
	if err != nil {
		return c.Status(http.StatusUnprocessableEntity).JSON(fiber.Map{
			"msg": http.StatusInternalServerError,
			"err": err.Error(),
		})
	}
	_, err = rec.Start()
	if err != nil {
		return c.Status(http.StatusUnprocessableEntity).JSON(fiber.Map{
			"msg": http.StatusInternalServerError,
			"err": err.Error(),
		})
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"code":    http.StatusOK,
		"message": "successful",
		"data": map[string]interface{}{
			"rid":     rec.RID,
			"sid":     rec.SID,
			"token":   rec.Token,
			"channel": rec.CallInfo.Channel,
			"uid":     rec.UID,
		},
	})
}

func stopRecording(c *fiber.Ctx) error {
	u := new(schemas.StopRecording)

	if err := c.BodyParser(u); err != nil {
		return c.Status(http.StatusUnprocessableEntity).JSON(fiber.Map{
			"msg": "invalid json",
			"err": err.Error(),
		})
	}

	data, err := utils.Stop(u.Channel, u.Uid, u.Rid, u.Sid)
	if err != nil {
		return c.Status(http.StatusUnprocessableEntity).JSON(fiber.Map{
			"msg": http.StatusInternalServerError,
			"err": err.Error(),
		})
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"code":    http.StatusOK,
		"message": "successful",
		"data":    data,
	})
}

func queryRecording(c *fiber.Ctx) error {
	u := new(schemas.QueryRecording)

	if err := c.BodyParser(u); err != nil {
		return c.Status(http.StatusUnprocessableEntity).JSON(fiber.Map{
			"msg": "invalid json",
			"err": err.Error(),
		})
	}

	data, err := utils.Query(u.Rid, u.Sid)
	if err != nil {
		return c.Status(http.StatusUnprocessableEntity).JSON(fiber.Map{
			"msg": http.StatusInternalServerError,
			"err": err.Error(),
		})
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"code":    http.StatusOK,
		"message": "successful",
		"data":    data,
	})
}

func updateRecording(c *fiber.Ctx) error {
	u := new(schemas.UpdateRecording)

	if err := c.BodyParser(u); err != nil {
		return c.Status(http.StatusUnprocessableEntity).JSON(fiber.Map{
			"msg": "invalid json",
			"err": err.Error(),
		})
	}

	_, err := utils.Update(u.Channel, u.Uid, u.Rid, u.Sid, u.Streamers)
	if err != nil {
		return c.Status(http.StatusUnprocessableEntity).JSON(fiber.Map{
			"msg": http.StatusInternalServerError,
			"err": err.Error(),
		})
	}

	_, err = utils.UpdateLayout(u.Channel, u.Uid, u.Rid, u.Sid, u.Streamers)
	if err != nil {
		return c.Status(http.StatusUnprocessableEntity).JSON(fiber.Map{
			"msg": http.StatusInternalServerError,
			"err": err.Error(),
		})
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"code":    http.StatusOK,
		"message": "successful",
	})
}

func createRTCToken(c *fiber.Ctx) error {
	channel := c.Params("channel")
	uid := int(rand.Uint32())
	rtcToken, err := utils.GetRtcToken(channel, uid)
	if err != nil {
		return c.Status(http.StatusUnprocessableEntity).JSON(fiber.Map{
			"msg": http.StatusInternalServerError,
			"err": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"code":      http.StatusOK,
		"rtc_token": rtcToken,
		"uid":       uid,
	})
}

func createRTMToken(c *fiber.Ctx) error {
	uid := c.Params("uid")
	rtmToken, err := utils.GetRtmToken(fmt.Sprint(uid))
	if err != nil {
		return c.Status(http.StatusUnprocessableEntity).JSON(fiber.Map{
			"msg": http.StatusInternalServerError,
			"err": err.Error(),
		})
	}
	return c.JSON(fiber.Map{
		"code":      http.StatusOK,
		"rtm_token": rtmToken,
	})
}

func createTokens(c *fiber.Ctx) error {
	channel := c.Params("channel")
	uid := int(rand.Uint32())
	rtcToken, err := utils.GetRtcToken(channel, uid)
	if err != nil {
		return c.Status(http.StatusUnprocessableEntity).JSON(fiber.Map{
			"msg": http.StatusInternalServerError,
			"err": err.Error(),
		})
	}
	rtmToken, err := utils.GetRtmToken(fmt.Sprint(uid))
	if err != nil {
		return c.Status(http.StatusUnprocessableEntity).JSON(fiber.Map{
			"msg": http.StatusInternalServerError,
			"err": err.Error(),
		})
	}
	return c.JSON(fiber.Map{
		"code":      http.StatusOK,
		"rtc_token": rtcToken,
		"rtm_token": rtmToken,
	})
}

// MountRoutes mounts all routes declared here
func MountRoutes(app *fiber.App) {
	app.Post("/api/recording/start", startRecording)
	app.Post("/api/recording/stop", stopRecording)
	app.Post("/api/recording/status", queryRecording)
	app.Post("/api/recording/update", updateRecording)
	app.Get("/api/get/rtc/:channel", createRTCToken)
	app.Get("/api/get/rtm/:uid", createRTMToken)
	app.Get("/api/tokens/:channel", createTokens)
}
