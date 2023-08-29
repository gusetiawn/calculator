package main

import (
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"time"
)

type LoanRequest struct {
	Plafon               float64 `json:"plafon"`
	LamaPinjaman         int     `json:"lama_pinjaman"`
	SukuBungaPertahun    float64 `json:"suku_bunga_pertahun"`
	TanggalMulaiAngsuran string  `json:"tanggal_mulai_angsuran"`
}

type Installment struct {
	AngsuranKe        int    `json:"angsuran_ke"`
	TanggalAngsuran   string `json:"tanggal_angsuran"`
	TotalAngsuran     int    `json:"total_angsuran"`
	AngsuranPokok     int    `json:"angsuran_pokok"`
	AngsuranBunga     int    `json:"angsuran_bunga"`
	SisaAngsuranPokok int    `json:"sisa_angsuran_pokok"`
}

func calculateInstallments(req LoanRequest) []Installment {
	var installments []Installment

	bungaBulanan := req.SukuBungaPertahun / 100 / 12
	tanggalAngsuran, _ := time.Parse("2006-01-02", req.TanggalMulaiAngsuran)

	sisaAngsuranPokok := int(req.Plafon)

	for i := 0; i < req.LamaPinjaman; i++ {
		angsuranKe := i + 1
		tanggalAngsuran = tanggalAngsuran.AddDate(0, 1, 0)
		totalAngsuran := calculateTotalAngsuran(req.Plafon, bungaBulanan, req.LamaPinjaman)

		angsuranBunga := calculateAngsuranBunga(sisaAngsuranPokok, bungaBulanan)
		angsuranPokok := totalAngsuran - angsuranBunga

		sisaAngsuranPokok -= angsuranPokok

		installment := Installment{
			AngsuranKe:        angsuranKe,
			TanggalAngsuran:   tanggalAngsuran.Format("2006-01-02"),
			TotalAngsuran:     totalAngsuran,
			AngsuranPokok:     angsuranPokok,
			AngsuranBunga:     angsuranBunga,
			SisaAngsuranPokok: sisaAngsuranPokok,
		}

		installments = append(installments, installment)
	}

	installments[len(installments)-1].SisaAngsuranPokok = 0

	return installments
}

func calculateTotalAngsuran(plafon, bungaBulanan float64, lamaPinjaman int) int {
	return int((plafon * bungaBulanan) / (1 - math.Pow(1+bungaBulanan, -float64(lamaPinjaman))))
}

func calculateAngsuranBunga(sisaAngsuranPokok int, bungaBulanan float64) int {
	return int(math.Round(float64(sisaAngsuranPokok) * bungaBulanan))
}

func calculateHandler(w http.ResponseWriter, r *http.Request) {
	var req LoanRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	installments := calculateInstallments(req)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(installments)
}

func main() {
	http.HandleFunc("/calculate", calculateHandler)
	fmt.Println("Server is running on port 8080")
	http.ListenAndServe(":8080", nil)
}
