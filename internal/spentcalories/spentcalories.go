package spentcalories

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"
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
	parts := strings.Split(data, ",")

	if len(parts) != 3 {
		return 0, "", 0, errors.New("неверный формат строки: ожидается три элемента, разделенных запятой (количество шагов, вид активности, продолжительность)")
	}

	steps, err := strconv.Atoi(parts[0])
	if err != nil {
		return 0, "", 0, errors.New("Ошибка преобразования количества шагов в число:" + err.Error())
	}

	if steps <= 0 {
		return 0, "", 0, errors.New("количество шагов должно быть больше нуля")
	}

	activity := parts[1]

	duration, err := time.ParseDuration(parts[2])
	if err != nil {
		return 0, "", 0, errors.New("ошибка парсинга продолжительности: " + err.Error())
	}

	return steps, activity, duration, nil // TODO: реализовать функцию
}

func distance(steps int, height float64) float64 {
	stepLength := height * stepLengthCoefficient

	distanceMeters := float64(steps) * stepLength

	distanceKm := distanceMeters / mInKm

	return distanceKm // TODO: реализовать функцию
}

func meanSpeed(steps int, height float64, duration time.Duration) float64 {
	if duration <= 0 {
		return 0
	}

	distanceKm := distance(steps, height)

	durationHours := duration.Hours()

	if durationHours == 0 {
		return 0
	}

	speed := distanceKm / durationHours

	return speed // TODO: реализовать функцию
}

func TrainingInfo(data string, weight, height float64) (string, error) {
	steps, activity, duration, err := parseTraining(data)
	if err != nil {
		log.Println(err)
		return "", err
	}

	distanceKm := distance(steps, height)

	speed := meanSpeed(steps, height, duration)

	var calories float64
	var errCalories error

	switch activity {
	case "Ходьба":
		calories, errCalories = WalkingSpentCalories(steps, weight, height, duration)
	case "Бег":
		calories, errCalories = RunningSpentCalories(steps, weight, height, duration)
	default:
		return "", errors.New("неизвестный тип тренировки")
	}
	if errCalories != nil {
		log.Println("Ошибка при рассчете калорий:", errCalories)
		return "", errCalories
	}

	result := fmt.Sprintf(
		"Тип тренировки: %s\nДлительность: %.2f ч.\nДистанция: %.2f км.\nСкорость: %.2f км/ч.\nСожгли калорий: %.2f",
		activity, duration.Hours(), distanceKm, speed, calories,
	)

	return result, nil
}

// TODO: реализовать функцию

func RunningSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	if steps <= 0 || weight <= 0 || height <= 0 || duration <= 0 {
		return 0, errors.New("некорректные входные параметры")
	}

	speed := meanSpeed(steps, height, duration)
	if speed == 0 {
		return 0, errors.New("не удалось рассчитать среднюю скорость (возможно, нулевая дистанция или продолжительность)")
	}

	durationMinutes := duration.Minutes()

	calories := (weight * speed * durationMinutes) / minInH

	return calories, nil // TODO: реализовать функцию
}

func WalkingSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	calories, err := RunningSpentCalories(steps, weight, height, duration)
	if err != nil {
		return 0, err
	}

	calories *= walkingCaloriesCoefficient
	return calories, nil // TODO: реализовать функцию
}
