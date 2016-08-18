package sakshat

import (
	"github.com/stianeikeland/go-rpio"
	"log"
	"os"
	"os/signal"
	"./entities"
)

var (
	Buzzer         *entities.Buzzer
	LEDRow         *entities.Led74HC595
	DigitalDisplay *entities.DigitalDisplayTM1637
	Ds18b20        *entities.DS18B20
)

func SaksGpioInit() {
	err := rpio.Open()
	if err != nil {
		log.Fatal(err)
	}

	process := []rpio.Pin{IC_TM1637_DI, IC_TM1637_CLK, IC_74HC595_DS, IC_74HC595_SHCP, IC_74HC595_STCP}
	for p := range (process) {
		process[p].Output()
		process[p].Low()
	}

	process = []rpio.Pin{BUZZER, TACT_RIGHT, TACT_LEFT, DIP_SWITCH_1, DIP_SWITCH_2}
	for p := range (process) {
		process[p].Output()
		process[p].High()
	}

	process = []rpio.Pin{TACT_RIGHT, TACT_LEFT, DIP_SWITCH_1, DIP_SWITCH_2}
	for p := range (process) {
		process[p].Input()
		process[p].PullUp()
	}
}

func init() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for _ = range c {
			log.Println("Closing pins and terminating program...")
			rpio.Close()
			os.Exit(0)
		}
	}()
	SaksGpioInit()
	Buzzer = &entities.Buzzer{
		Pin:      BUZZER,
		RealTrue: rpio.Low,
	}
	LEDRow = &entities.Led74HC595{
		IC: &entities.IC_74HC595{
			Pins:     map[string]rpio.Pin{"ds": IC_74HC595_DS, "shcp": IC_74HC595_SHCP, "stcp": IC_74HC595_STCP},
			RealTrue: rpio.High,
		},
	}
	Ds18b20 = &entities.DS18B20{
		Pin: DS18B20,
	}
	DigitalDisplay = &entities.DigitalDisplayTM1637{
		IC: &entities.IC_TM1637{
			Pins: map[string]rpio.Pin{"di": IC_TM1637_DI, "clk": IC_TM1637_CLK},
			RealTrue: rpio.High,
		},
	}
}