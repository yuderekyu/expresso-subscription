package handlers

import (
	"github.com/pborman/uuid"
	"gopkg.in/alexcesaro/statsd.v2"
	"gopkg.in/gin-gonic/gin.v1"

	"github.com/ghmeier/bloodlines/handlers"
	"github.com/yuderekyu/covenant/helpers"
	"github.com/yuderekyu/covenant/models"
)

/*SubscriptionI contains the methods for a subscription handler*/
type SubscriptionI interface {
	New(ctx *gin.Context)
	ViewAll(ctx *gin.Context)
	View(ctx *gin.Context)
	ViewByRoaster(ctx *gin.Context)
	ViewByUser(ctx *gin.Context)
	Update(ctx *gin.Context)
	Delete(ctx *gin.Context)
	Time() gin.HandlerFunc
}

/*Subscription is the handler for all subscription api calls*/
type Subscription struct {
	*handlers.BaseHandler
	Subscription helpers.SubscriptionI
}

/*NewSubscription returns a subscription handler*/
func NewSubscription(ctx *handlers.GatewayContext) SubscriptionI {
	stats := ctx.Stats.Clone(statsd.Prefix("api.subscription"))
	return &Subscription{
		Subscription:      helpers.NewSubscription(ctx.Sql, ctx.TownCenter, ctx.Warehouse),	
		BaseHandler: &handlers.BaseHandler{Stats: stats}, 
	}
}

/*New adds the given subscription entry to the database*/
func (s *Subscription) New(ctx *gin.Context) {

	var json models.RequestIdentifiers
	err := ctx.BindJSON(&json)
	if err != nil {
		s.UserError(ctx, "Error: unable to parse json", err)
		return
	}

	subscription := models.NewSubscription(json.UserID, json.Frequency, json.RoasterID, json.ItemID)
	err = s.Subscription.Insert(subscription)
	if err != nil {
		s.ServerError(ctx, err, json)
		return
	}

	s.Success(ctx, subscription)
}

/*View returns the subscription entry with the given subscriptionId*/
func (s *Subscription) View(ctx *gin.Context) {
	id := ctx.Param("subscriptionId")

	subscription, err := s.Subscription.GetByID(id)
	if err != nil {
		s.ServerError(ctx, err, nil)
		return
	}

	if(subscription == nil) {
		s.UserError(ctx, "Error: Subscription does not exist", id)
		return
	}

	s.Success(ctx, subscription)
}

/*ViewAll returns a list of subscriptions with offset and limit determining the entries and amount*/
func (s *Subscription) ViewAll(ctx *gin.Context) {
	offset, limit := s.GetPaging(ctx)
	subscriptions, err := s.Subscription.GetAll(offset, limit)
	if err != nil {
		s.ServerError(ctx, err, nil)
		return
	}

	s.Success(ctx, subscriptions)
}

/*ViewByRoaster returns a list of subscriptions with the given roasterId,
with the offset and limit determining the entries and amount*/
func (s *Subscription) ViewByRoaster(ctx *gin.Context) {
	roasterID := ctx.Param("roasterId")
	offset, limit := s.GetPaging(ctx)

	subscriptions, err := s.Subscription.GetByRoaster(roasterID, offset, limit)
	if err != nil {
		s.ServerError(ctx, err, nil)
		return
	}

	s.Success(ctx, subscriptions)
}

/*ViewByUser returns a list of subscriptions with the given roasterId,
with the offset and limit determining the entries and amount*/
func (s *Subscription) ViewByUser(ctx *gin.Context) {
	userID := ctx.Param("userId")
	offset, limit := s.GetPaging(ctx)

	subscriptions, err := s.Subscription.GetByUser(userID, offset, limit)
	if err != nil {
		s.ServerError(ctx, err, nil)
		return
	}

	s.Success(ctx, subscriptions)
}

/*Update overwrites a subscription with the given subscriptionId*/
func (s *Subscription) Update(ctx *gin.Context) {
	id := ctx.Param("subscriptionId")

	var json models.Subscription
	err := ctx.BindJSON(&json)
	if err != nil {
		s.UserError(ctx, "Error: unable to parse json", err)
		return
	}
	json.ID = uuid.Parse(id)

	err = s.Subscription.Update(id, &json)
	if err != nil {
		s.ServerError(ctx, err, json)
		return
	}

	s.Success(ctx, json)
}

/*Delete removes a subscription with the given subscriptionId*/
func (s *Subscription) Delete(ctx *gin.Context) {
	id := ctx.Param("subscriptionId")

	err := s.Subscription.Delete(id)
	if err != nil {
		s.ServerError(ctx, err, id)
		return
	}
	s.Success(ctx, nil)
}
