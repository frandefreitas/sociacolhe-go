package services

import "strings"

// AutoTriage: heurística simples para MVP.
// Produção: substituir por regra mais robusta/modelo.
func AutoTriage(desc string) (area string, urgency int) {
	d := strings.ToLower(desc)
	urgency = 2
	if strings.Contains(d, "suicid") || strings.Contains(d, "autoagress") || strings.Contains(d, "autoles") {
		urgency = 5
	}
	if strings.Contains(d, "ansiedad") || strings.Contains(d, "ansiedade") { area = "Ansiedade" }
	if strings.Contains(d, "depress") { area = "Depressão" }
	if strings.Contains(d, "violenc") { area = "Violência" }
	if area == "" { area = "Geral" }
	return
}
