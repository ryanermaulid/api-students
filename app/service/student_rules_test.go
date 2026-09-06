package service

import (
	"testing"

	"api-students/app/model"
)

func TestValidateCreate(t *testing.T) {
	cases := []struct {
		name     string
		req      model.CreateStudentRequest
		wantErrs int
	}{
		{"Valid Request", model.CreateStudentRequest{NIM: "123", Name: "Ryan"}, 0},
		{"Empty NIM", model.CreateStudentRequest{NIM: "  ", Name: "Ryan"}, 1},
		{"Empty All", model.CreateStudentRequest{NIM: "", Name: ""}, 2},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			errs := ValidateCreate(tc.req)
			if len(errs) != tc.wantErrs {
				t.Errorf("harap %d error, dapat %d", tc.wantErrs, len(errs))
			}
		})
	}
}

func TestApplyPatch(t *testing.T) {
	initial := model.Student{ID: 1, NIM: "101", Name: "Lama", Grade: 80, IsActive: true}
	newName := "Baru"
	
	req := model.PatchStudentRequest{Name: &newName}
	result, errs := ApplyPatch(initial, req)

	if len(errs) != 0 {
		t.Fatalf("tidak seharusnya ada error, dapat: %v", errs)
	}
	if result.Name != "Baru" {
		t.Errorf("harap nama berubah menjadi Baru, dapat %s", result.Name)
	}
	if result.NIM != "101" {
		t.Error("field NIM yang tidak dikirim seharusnya tidak berubah")
	}
}

func TestCountTotalPages(t *testing.T) {
	cases := []struct{ total, limit, want int }{
		{0, 10, 0},
		{1, 10, 1},
		{10, 10, 1},
		{11, 10, 2},
		{137, 20, 7},
	}

	for _, tc := range cases {
		if got := CountTotalPages(tc.total, tc.limit); got != tc.want {
			t.Errorf("total=%d limit=%d: harap %d, dapat %d", tc.total, tc.limit, tc.want, got)
		}
	}
}