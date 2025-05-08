package daysteps

import (
	"time"
	"strings"
	"errors"
	"strconv"
	"fmt"
	"log"

	spentCalories "github.com/Yandex-Practicum/tracker/internal/spentcalories"
)

const (
	// Длина одного шага в метрах
	stepLength = 0.65
	// Количество метров в одном километре
	mInKm = 1000
)

func parsePackage(data string) (int, time.Duration, error) {

	splited := strings.Split(data, ",")
	if len(splited) != 2{
		return 0, time.Duration(0), errors.New("не корректная строка ")
	}

	steps, err := strconv.Atoi(splited[0])
	if err != nil{
		return 0, time.Duration(0), err
	}
	if steps <= 0{
		return 0, time.Duration(0), errors.New("не корректное кол-во шагов ")
	}

	walkDuration, err := time.ParseDuration(splited[1])
	if err != nil{
		return 0, time.Duration(0), err
	}
	if walkDuration <= 0{
		return 0, time.Duration(0), errors.New("не корректная длительность ")
	}

	return steps, walkDuration, nil
} 

func DayActionInfo(data string, weight, height float64) string {

	steps, walkDuration, err := parsePackage(data)
	if err != nil{
		log.Printf("Ошибка: %v", err)
		return ""
	}
	if steps <= 0{
		return ""
	}

	distanceInM := float64(steps) * stepLength
	distanceInKm := distanceInM / mInKm

	calories, err := spentCalories.WalkingSpentCalories(steps, weight, height, walkDuration)
	if err != nil{
		log.Printf("Ошибка: %v", err)
		return ""
	}

	return fmt.Sprintf("Количество шагов: %d.\nДистанция составила %.2f км.\nВы сожгли %.2f ккал.\n", steps, distanceInKm, calories)
}
