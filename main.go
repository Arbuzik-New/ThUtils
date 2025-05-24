package main

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"log"
	"net/http"
	"slices"
	"time"
)

var lastTime time.Time

type NickNamesResponse struct {
	Nicknames []string `json:"nicknames"`
}

type CoordinatesRequest struct {
	Nickname    string `json:"nickname"`
	Coordinates string `json:"coordinates"`
}

type CoordinatesResponse struct {
	Nickname    string `json:"nickname"`
	Coordinates string `json:"coordinates"`
}

type Response struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

var nicknames []string = make([]string, 0)

var coordinates string
var nickname string

func init() {
	lastTime = time.Now()
}

func clearNicknames() {
	for {
		if time.Now().Sub(lastTime).Minutes() > 5 {
			lastTime = time.Now()
			nicknames = make([]string, 0)
		}
	}
}

func main() {

	go clearNicknames()

	var echo *echo.Echo = echo.New()

	// Coordinates
	echo.GET("/coordinates", GetCoordinates)
	echo.POST("/coordinates", PostCoordinates)
	echo.POST("/clear/coordinates", PostClearCoordinates)

	// Nicknames
	echo.GET("/nicknames", GetNicknames)
	echo.POST("/nicknames", PostNicknames)

	var err error = echo.Start("localhost:8080")
	if err != nil {
		log.Fatal(err)
	}
}

func GetCoordinates(context echo.Context) error {
	return context.JSON(http.StatusOK, &CoordinatesResponse{nickname, coordinates})
}

func GetNicknames(context echo.Context) error {
	return context.JSON(http.StatusOK, &NickNamesResponse{nicknames})
}

func PostCoordinates(context echo.Context) error {
	var coords CoordinatesRequest
	var err error = context.Bind(&coords)

	if err != nil {
		fmt.Println(err.Error())
		return context.JSON(http.StatusBadRequest, &Response{"Bad Request", "Could set the coordinates"})
	}

	nickname = coords.Nickname
	coordinates = coords.Coordinates

	return context.JSON(http.StatusOK, &Response{"Success", "Set the coordinates"})
}

func PostClearCoordinates(context echo.Context) error {
	nickname = ""
	coordinates = ""

	return context.JSON(http.StatusOK, &Response{"Success", "Cleared the coordinates"})
}

func PostNicknames(context echo.Context) error {
	var nickname string
	var err error = context.Bind(&nickname)
	if err != nil {
		fmt.Println(err.Error())
		return context.JSON(http.StatusBadRequest, &Response{"Bad Request", "Could add the nickname"})
	}
	if !slices.Contains(nicknames, nickname) {
		nicknames = append(nicknames, nickname)
	}
	return context.JSON(http.StatusOK, &Response{"Success", "Added the nickname"})
}
