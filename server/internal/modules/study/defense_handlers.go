package study

import (
	"net/http"

	"github.com/babakapaa/lodestar/server/internal/platform/jsonhttp"
)

type SlotsHandler struct{ service DefenseService }

func (h SlotsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	list, err := h.service.Slots(r.Context(), principal(r), r.PathValue("id"))
	respond(w, http.StatusOK, list, err)
}

type AddSlotsHandler struct{ service DefenseService }

func (h AddSlotsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var in SlotsInput
	if err := jsonhttp.Read(w, r, &in); err != nil {
		badRequest(w)
		return
	}
	list, err := h.service.AddSlots(r.Context(), principal(r), r.PathValue("id"), in)
	respond(w, http.StatusCreated, list, err)
}

type DeleteSlotHandler struct{ service DefenseService }

func (h DeleteSlotHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if err := h.service.DeleteSlot(r.Context(), principal(r), r.PathValue("id")); err != nil {
		writeError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

type MyDefenseHandler struct{ service DefenseService }

func (h MyDefenseHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	my, err := h.service.Mine(r.Context(), principal(r), r.PathValue("id"))
	respond(w, http.StatusOK, my, err)
}

type BookHandler struct{ service DefenseService }

func (h BookHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var in BookingInput
	if err := jsonhttp.Read(w, r, &in); err != nil {
		badRequest(w)
		return
	}
	my, err := h.service.Book(r.Context(), principal(r), r.PathValue("id"), in)
	respond(w, http.StatusOK, my, err)
}

type UnbookHandler struct{ service DefenseService }

func (h UnbookHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	my, err := h.service.Unbook(r.Context(), principal(r), r.PathValue("id"))
	respond(w, http.StatusOK, my, err)
}

type DefensesHandler struct{ service DefenseService }

func (h DefensesHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	list, err := h.service.List(r.Context(), principal(r), r.PathValue("id"))
	respond(w, http.StatusOK, list, err)
}

type SampleHandler struct{ service DefenseService }

func (h SampleHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var in SampleInput
	if err := jsonhttp.Read(w, r, &in); err != nil {
		badRequest(w)
		return
	}
	chosen, err := h.service.Sample(r.Context(), principal(r), r.PathValue("id"), in)
	respond(w, http.StatusOK, map[string]any{"studentIds": chosen}, err)
}

type DefenseHandler struct{ service DefenseService }

func (h DefenseHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	d, err := h.service.Get(r.Context(), principal(r), r.PathValue("id"), r.PathValue("student"))
	respond(w, http.StatusOK, d, err)
}

type SaveDefenseHandler struct{ service DefenseService }

func (h SaveDefenseHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var in DefenseInput
	if err := jsonhttp.Read(w, r, &in); err != nil {
		badRequest(w)
		return
	}
	d, err := h.service.Save(r.Context(), principal(r), r.PathValue("id"), r.PathValue("student"), in)
	respond(w, http.StatusOK, d, err)
}

type SelectHandler struct{ service DefenseService }

func (h SelectHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var in SelectionInput
	if err := jsonhttp.Read(w, r, &in); err != nil {
		badRequest(w)
		return
	}
	d, err := h.service.Select(r.Context(), principal(r), r.PathValue("id"), r.PathValue("student"), in)
	respond(w, http.StatusOK, d, err)
}

type FinishHandler struct{ service DefenseService }

func (h FinishHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var in FinishInput
	if err := jsonhttp.Read(w, r, &in); err != nil {
		badRequest(w)
		return
	}
	d, err := h.service.Finish(r.Context(), principal(r), r.PathValue("id"), r.PathValue("student"), in)
	respond(w, http.StatusOK, d, err)
}

func respond(w http.ResponseWriter, status int, body any, err error) {
	if err != nil {
		writeError(w, err)
		return
	}
	jsonhttp.Write(w, status, body)
}
