package handlers

import (
	"dynamic-user-segmentation-service/internal/errors"
	"dynamic-user-segmentation-service/internal/repositories"
	"encoding/json"
	"net/http"

	"dynamic-user-segmentation-service/internal/models"
	"github.com/gorilla/mux"
)

func GetSegmentsHandler(w http.ResponseWriter, r *http.Request) {
	segmentRepository := repositories.PostgresSegmentRepository{}

	segments, err := segmentRepository.GetAllSegments()

	w.Header().Add("Content-Type", "application/json")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		errorDto := &models.ErrorDto{
			Error: "Возникла внутренняя ошибка при запросе всех сегментов",
		}
		json.NewEncoder(w).Encode(errorDto)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(segments)
}

func GetSegmentBySlugHandler(w http.ResponseWriter, r *http.Request) {
	segmentRepository := repositories.PostgresSegmentRepository{}

	params := mux.Vars(r)
	slug := params["slug"]

	segment, err := segmentRepository.GetSegmentBySlug(slug)
	w.Header().Add("Content-Type", "application/json")
	if err != nil {
		switch err {
		case errors.RecordNotFound:
			w.WriteHeader(http.StatusNotFound)
			errorDto := &models.ErrorDto{
				Error: "Запись с таким названием в таблице сегментов не найдена",
			}
			json.NewEncoder(w).Encode(errorDto)
		default:
			w.WriteHeader(http.StatusInternalServerError)
			errorDto := &models.ErrorDto{
				Error: "Возникла внутренняя ошибка при запросе сегмента по названию",
			}
			json.NewEncoder(w).Encode(errorDto)
		}
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(segment)
}

func CreateSegmentHandler(w http.ResponseWriter, r *http.Request) {
	segmentRepository := repositories.PostgresSegmentRepository{}

	var segment models.CreateOrUpdateSegmentDto

	err := json.NewDecoder(r.Body).Decode(&segment)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		errorDto := &models.ErrorDto{
			Error: "Некорректные входные данные",
		}
		json.NewEncoder(w).Encode(errorDto)
		return
	}

	err = segmentRepository.CreateSegment(segment.Slug, segment.Description)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		errorDto := &models.ErrorDto{
			Error: "Возникла внутренняя ошибка при создании сегмента",
		}
		w.Header().Add("Content-Type", "application/json")
		json.NewEncoder(w).Encode(errorDto)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func UpdateSegmentHandler(w http.ResponseWriter, r *http.Request) {
	segmentRepository := repositories.PostgresSegmentRepository{}

	params := mux.Vars(r)
	slug := params["slug"]

	var segment models.CreateOrUpdateSegmentDto

	err := json.NewDecoder(r.Body).Decode(&segment)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		errorDto := &models.ErrorDto{
			Error: "Некорректные входные данные",
		}
		json.NewEncoder(w).Encode(errorDto)
		return
	}

	err = segmentRepository.UpdateSegment(slug, segment.Description)
	if err != nil {
		w.Header().Add("Content-Type", "application/json")
		switch err {
		case errors.RecordNotFound:
			w.WriteHeader(http.StatusNotFound)
			errorDto := &models.ErrorDto{
				Error: "Запись с таким названием в таблице сегментов не найдена",
			}
			json.NewEncoder(w).Encode(errorDto)
		default:
			w.WriteHeader(http.StatusInternalServerError)
			errorDto := &models.ErrorDto{
				Error: "Возникла внутренняя ошибка при обновлении сегмента",
			}
			json.NewEncoder(w).Encode(errorDto)
		}
		return
	}

	w.WriteHeader(http.StatusOK)
}

func DeleteSegmentHandler(w http.ResponseWriter, r *http.Request) {
	segmentRepository := repositories.PostgresSegmentRepository{}

	params := mux.Vars(r)
	slug := params["slug"]

	_, err := segmentRepository.DeleteSegment(slug)
	if err != nil {
		w.Header().Add("Content-Type", "application/json")
		switch err {
		case errors.RecordNotFound:
			w.WriteHeader(http.StatusNotFound)
			errorDto := &models.ErrorDto{
				Error: "Запись с таким названием в таблице сегментов не найдена",
			}
			json.NewEncoder(w).Encode(errorDto)
		default:
			w.WriteHeader(http.StatusInternalServerError)
			errorDto := &models.ErrorDto{
				Error: "Возникла внутренняя ошибка при удалении сегмента",
			}
			json.NewEncoder(w).Encode(errorDto)
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
