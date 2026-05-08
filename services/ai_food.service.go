package services

import (
	"regexp"
	"strconv"
	"strings"

	"github.com/vpa/quanlynhahang-backend/config"
	"github.com/vpa/quanlynhahang-backend/models"
)

var userState = map[string]string{}

func HandleFoodChat(userID, message string) models.AIChatResponse {

	msg := strings.ToLower(strings.TrimSpace(message))

	// 1. Greeting
	if res, ok := handleGreeting(userID, msg); ok {
		return res
	}

	// 2. Đang chờ nhập giá
	if userState[userID] == "WAIT_PRICE" {
		return handleWaitingPrice(userID, msg)
	}

	// 3. Hỏi giá
	if res, ok := handleAskPrice(userID, msg); ok {
		return res
	}

	// 4. Gợi ý chung
	return handleSuggestion(msg)
}

// ===== GREETING =====
func handleGreeting(userID, msg string) (models.AIChatResponse, bool) {

	keywordsSayHelo := []string{"hi", "hello", "xin chào", "chào", "hey"}
	keywordsSayBye := []string{"bye","BYEBYE","BYE","goodbye","tam biet","TAMBIET","Tạm Biệt"}

	for _, k := range keywordsSayHelo {
		if strings.Contains(msg, k) {
			delete(userState, userID)
			return text("Xin chào bạn 👋 Nhà hàng có thể giúp gì cho bạn?"), true
		}
	}
	for _, k := range keywordsSayBye {
		if strings.Contains(msg, k) {
			delete(userState, userID)
			return text("Xin chào và hẹn gặp lại bạn. Chúc bạn và gia đình có 1 bữa tối tuyệt vời!!!"), true
		}
	}

	return models.AIChatResponse{}, false
}

// ===== ASK PRICE =====
func handleAskPrice(userID, msg string) (models.AIChatResponse, bool) {
	if strings.Contains(msg, "rẻ") || strings.Contains(msg, "giá") {
		userState[userID] = "WAIT_PRICE"
		return text("Bạn muốn món ăn giá khoảng bao nhiêu? (VD: 30000, 50k)"), true
	}
	return models.AIChatResponse{}, false
}

// ===== WAIT PRICE =====
func handleWaitingPrice(userID, msg string) models.AIChatResponse {

	price, ok := extractPriceFromMessage(msg)
	if !ok {
		return text("Bạn có thể nhập giá như: 50k, 80000, dưới 100k")
	}

	var monAn []models.MonAn
	config.DB.
		Preload("AnhMonAn").
		Where("gia_tien <= ?", price).
		Order("gia_tien DESC").
		Limit(5).
		Find(&monAn)

	// ❗ KHÔNG reset state nếu chưa có món
	if len(monAn) == 0 {
		return text("Chưa có món trong tầm giá này 😢 Bạn thử tăng giá lên nhé")
	}

	// ✅ CÓ MÓN → reset state
	delete(userState, userID)

	var foods []models.FoodChat
	for _, m := range monAn {
		img := ""
		if len(m.AnhMonAn) > 0 {
			img = m.AnhMonAn[0].ImageURL
		}

		foods = append(foods, models.FoodChat{
			MaMonAn: int(m.MaMonAn),
			Ten:     m.TenMonAn,
			Gia:     m.GiaTien,
			Anh:     img,
		})
	}

	return models.AIChatResponse{
		Type:    "food_list",
		Message: "Đây là các món phù hợp với mức giá bạn nhập:",
		Foods:   foods,
	}
}

// ===== SUGGEST =====
func handleSuggestion(msg string) models.AIChatResponse {
	switch {
	case strings.Contains(msg, "cay"):
		return text("Món cay gợi ý: Mì cay, Gà sốt cay, Lẩu Thái")
	case strings.Contains(msg, "ngọt"):
		return text("Món ngọt: Bánh flan, Chè xoài, Trà sữa")
	case strings.Contains(msg, "2 người"):
		return text("Combo cho 2 người: Lẩu nhỏ + nước uống")
	case strings.Contains(msg, "ăn gì"):
		return text("Hôm nay bạn có thể thử: Cơm gà, Bún bò, Mì Ý")
	default:
		return text("Chúng tôi có thể giúp gì cho bạn???")
	}
}

// ===== HELPERS =====
func text(msg string) models.AIChatResponse {
	return models.AIChatResponse{
		Type:    "text",
		Message: msg,
	}
}
// Đồi giá tiền từ chữ ra số
func extractPriceFromMessage(msg string) (int, bool) {

	msg = strings.ToLower(msg)
	msg = strings.ReplaceAll(msg, ".", "")
	msg = strings.ReplaceAll(msg, ",", "")
	msg = strings.ReplaceAll(msg, "k", "000")
	msg = strings.ReplaceAll(msg, "nghìn", "000")

	re := regexp.MustCompile(`\d+`)
	match := re.FindString(msg)

	if match == "" {
		return 0, false
	}

	price, err := strconv.Atoi(match)
	if err != nil {
		return 0, false
	}

	return price, true
}