package spentcalories

import (
	"time"
	"errors"
	"strconv"
	"strings"
	"fmt"
)

// Основные константы, необходимые для расчетов.
const (
	lenStep                    = 0.65 // средняя длина шага.
	mInKm                      = 1000 // количество метров в километре.
	minInH                     = 60   // количество минут в часе.
	stepLengthCoefficient      = 0.45 // коэффициент для расчета длины шага на основе роста.
	walkingCaloriesCoefficient = 0.5  // коэффициент для расчета калорий при ходьбе
)

func parseTraining(data string) (int, string, time.Duration, error) {

	splited := strings.Split(data, ",")
	if len(splited) != 3{
		return 0, "", time.Duration(0), errors.New("не корректная строка")
	}
	
	steps, err := strconv.Atoi(splited[0])
	if err != nil{
		return 0, "", time.Duration(0), err
	}
	if steps <= 0{
		return 0, "", time.Duration(0), errors.New("не корректное кол-во шагов")
	}

	activityDuration, err := time.ParseDuration(splited[2])
	if err != nil{
		return 0, "", time.Duration(0), err
	}
	if activityDuration <= 0{
		return 0, "", time.Duration(0), errors.New("не корректная длительность")
	}

	return steps, splited[1], activityDuration, nil
}

func distance(steps int, height float64) float64 {

	stepLength := height * stepLengthCoefficient

	distanceInKm := float64(steps) * stepLength / mInKm

	return distanceInKm
}

func meanSpeed(steps int, height float64, duration time.Duration) float64 {
	
	if duration <= 0{
		return 0
	}

	durationInHours := duration.Hours()

	return  distance(steps, height) / durationInHours
}

func TrainingInfo(data string, weight, height float64) (string, error) {
	
	steps, activityType, activityDuration, err := parseTraining(data)
	if err != nil{
		return "", err
	}

	switch activityType{
	case "Ходьба":
		calories, err := WalkingSpentCalories(steps, weight, height, activityDuration)
		if err != nil{
			return "", err
		}

		speed := meanSpeed(steps, height, activityDuration)
		distance := distance(steps, height)

		return fmt.Sprintf("Тип тренировки: %s\nДлительность: %.2f ч.\nДистанция: %.2f км.\nСкорость: %.2f км/ч\nСожгли калорий: %.2f\n", activityType, activityDuration.Hours(), distance, speed, calories), nil
	
	case "Бег":
		calories, err := RunningSpentCalories(steps, weight, height, activityDuration)
		if err != nil{
			return "", err
		}

		speed := meanSpeed(steps, height, activityDuration)
		distance := distance(steps, height)

		return fmt.Sprintf("Тип тренировки: %s\nДлительность: %.2f ч.\nДистанция: %.2f км.\nСкорость: %.2f км/ч\nСожгли калорий: %.2f\n", activityType, activityDuration.Hours(), distance, speed, calories), nil
	default:
		return "", errors.New("неизвестный тип тренировки")
	}
	
}

func RunningSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	if steps <= 0 || weight <= 0 || height <= 0 || duration <= 0{
		return 0, errors.New("неккоректные параметры")
	}

	speed := meanSpeed(steps, height, duration)

	durationInMinutes := duration.Minutes()

	return (weight * speed * durationInMinutes) / minInH, nil
}

func WalkingSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	if steps <= 0 || weight <= 0 || height <= 0 || duration <= 0{
		return 0, errors.New("неккоректные параметры")
	}

	speed := meanSpeed(steps, height, duration)

	durationInMinutes := duration.Minutes()

	return (weight * speed * durationInMinutes) / minInH * walkingCaloriesCoefficient, nil
}
