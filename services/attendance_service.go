package services

import (
	"LecCenterGo/dao"
	"LecCenterGo/models"
	"errors"
	"gorm.io/gorm"
	"time"
)

// 创建打卡记录
func CreateAttendance(attendance *models.Attendance) error {
	// 检查今日是否已打卡（可选逻辑）
	var count int64
	today := time.Now().Format("2006-01-02")
	err := dao.DB.Model(&models.Attendance{}).
		Where("email = ? AND DATE(created_at) = ?", attendance.Email, today).
		Count(&count).Error

	if err != nil {
		return err
	}

	if count > 0 {
		return errors.New("今日已打卡")
	}

	return dao.DB.Create(attendance).Error
}

// 获取打卡记录列表（分页）
func GetAttendances(page, pageSize int) ([]models.Attendance, int64, error) {
	var attendances []models.Attendance
	var total int64

	// 查询总数
	if err := dao.DB.Model(&models.Attendance{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 查询分页数据（按时间倒序）
	if err := dao.DB.
		Order("created_at DESC").
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Find(&attendances).Error; err != nil {
		return nil, 0, err
	}

	return attendances, total, nil
}

// 按日期获取打卡记录
func GetAttendancesByDate(date string, page, pageSize int) ([]models.Attendance, int64, error) {
	var attendances []models.Attendance
	var total int64

	// 查询总数
	if err := dao.DB.Model(&models.Attendance{}).
		Where("DATE(created_at) = ?", date).
		Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 查询分页数据
	if err := dao.DB.
		Where("DATE(created_at) = ?", date).
		Order("created_at DESC").
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Find(&attendances).Error; err != nil {
		return nil, 0, err
	}

	return attendances, total, nil
}

// 获取单个打卡记录
func GetAttendanceByID(id uint) (*models.Attendance, error) {
	var attendance models.Attendance
	if err := dao.DB.First(&attendance, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &attendance, nil
}

// 导出指定日期的打卡记录
func ExportAttendancesByDate(date string) ([]models.Attendance, error) {
	var attendances []models.Attendance
	if err := dao.DB.
		Where("DATE(created_at) = ?", date).
		Order("created_at ASC").
		Find(&attendances).Error; err != nil {
		return nil, err
	}
	return attendances, nil
}
