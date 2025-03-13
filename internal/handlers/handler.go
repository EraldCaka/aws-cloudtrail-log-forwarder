package handlers

import "github.com/EraldCaka/aws-cloudtrail-log-forwarder/internal/services"

type Handler struct {
	redis services.RedisService
	aws   services.AWSService
	mongo services.MongoService
}
