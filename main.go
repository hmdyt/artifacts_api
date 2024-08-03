//go:generate go run github.com/ogen-go/ogen/cmd/ogen --target ./openapi -package openapi --clean ope

package main

import (
	"context"
	"log"
	"os"
	"sync"
	"time"

	"artifacts/openapi"

	"github.com/go-faster/errors"
)

func getToken() string {
	return os.Getenv("ARTIFACTS_TOKEN")
}

type SecuritySource struct {
}

func (s SecuritySource) HTTPBasic(_ context.Context, _ string) (openapi.HTTPBasic, error) {
	return openapi.HTTPBasic{}, nil
}

func (s SecuritySource) JWTBearer(_ context.Context, _ string) (openapi.JWTBearer, error) {
	return openapi.JWTBearer{Token: getToken()}, nil
}

var (
	client *openapi.Client
	ctx    context.Context
)

func init() {
	initClient()
	ctx = context.Background()
}

func initClient() {
	c, err := openapi.NewClient("https://api.artifactsmmo.com", SecuritySource{})
	if err != nil {
		log.Fatalf("Client init error: %v", err)
	}
	client = c
}

func main() {
	var wg sync.WaitGroup
	wg.Add(1)

	go loopAction("chami", 0, 1, Fight)
	go loopAction("moe", 2, 0, Gathering)
	go loopAction("ginji", -1, 0, Gathering)

	wg.Wait()
}

func loopAction(name string, x, y int, action func(string) error) {
	// TODO: 今がwait状態かどうか調べて必要なら待つ処理

	if err := MaybeMove(name, x, y); err != nil {
		log.Fatalf("MaybeMove Error: %v", err)
	}

	for {
		if err := action(name); err != nil {
			log.Fatalf("Error: %v", err)
		}
	}
}

func MaybeMove(name string, x, y int) error {
	res, err := client.ActionMoveMyNameActionMovePost(
		ctx,
		&openapi.DestinationSchema{
			X: x,
			Y: y,
		},
		openapi.ActionMoveMyNameActionMovePostParams{Name: name},
	)
	if err != nil {
		return err
	}

	switch r := res.(type) {
	case *openapi.CharacterMovementResponseSchema:
		wait(r.Data.Cooldown.Expiration, name, string(r.Data.Cooldown.Reason))
		return nil
	case *openapi.ActionMoveMyNameActionMovePostCode490:
		return nil
	}

	return errors.New("invalid response type")
}

func Gathering(name string) error {
	res, err := client.ActionGatheringMyNameActionGatheringPost(ctx,
		openapi.ActionGatheringMyNameActionGatheringPostParams{Name: name})
	if err != nil {
		return err
	}

	r, ok := res.(*openapi.SkillResponseSchema)
	if !ok {
		return err
	}

	wait(r.Data.Cooldown.Expiration, name, string(r.Data.Cooldown.Reason))
	return nil
}

func Fight(name string) error {
	res, err := client.ActionFightMyNameActionFightPost(ctx,
		openapi.ActionFightMyNameActionFightPostParams{Name: name})
	if err != nil {
		return err
	}

	r, ok := res.(*openapi.CharacterFightResponseSchema)
	if !ok {
		return err
	}

	wait(r.Data.Cooldown.Expiration, name, string(r.Data.Cooldown.Reason))
	return nil
}

func wait(until time.Time, name string, reason string) {
	diff := until.Sub(time.Now())
	go func() {
		log.Printf("[%s:%s] Waiting for %v sec", name, reason, diff.Seconds())
	}()
	time.Sleep(diff)
}
