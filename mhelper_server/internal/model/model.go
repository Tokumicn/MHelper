package model

import (
	"context"
	"fmt"
	"runtime"
	"time"

	otgorm "github.com/eddycjy/opentracing-gorm"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"mhelper_server/global"
	"mhelper_server/pkg/setting"
)

const (
	STATE_OPEN            = 1
	STATE_CLOSE           = 0
	SlowQueryStartTimeTag = "q_start_time"
)

type Model struct {
	ID         uint32 `gorm:"primary_key" json:"id"`
	CreatedBy  string `json:"created_by"`
	ModifiedBy string `json:"modified_by"`
	CreatedOn  uint32 `json:"created_on"`
	ModifiedOn uint32 `json:"modified_on"`
	DeletedOn  uint32 `json:"deleted_on"`
	IsDel      uint8  `json:"is_del"`
}

func NewDBEngine(databaseSetting *setting.DatabaseSettingS) (*gorm.DB, error) {
	s := "%s:%s@tcp(%s)/%s?charset=%s&parseTime=%t&loc=Local"
	db, err := gorm.Open(databaseSetting.DBType, fmt.Sprintf(s,
		databaseSetting.UserName,
		databaseSetting.Password,
		databaseSetting.Host,
		databaseSetting.DBName,
		databaseSetting.Charset,
		databaseSetting.ParseTime,
	))
	if err != nil {
		return nil, err
	}

	if global.ServerSetting.RunMode == "debug" {
		db.LogMode(true)
	}
	db.SingularTable(true)
	db.Callback().Create().Replace("gorm:update_time_stamp", updateTimeStampForCreateCallback)
	db.Callback().Update().Replace("gorm:update_time_stamp", updateTimeStampForUpdateCallback)
	db.Callback().Delete().Replace("gorm:delete", deleteCallback)

	// 慢SQL日志输出
	db.Callback().Query().Before("gorm:query").Register("slow:before", beforeQuery)
	db.Callback().Query().Before("gorm:after_query").Register("slow:after", afterQuery)

	db.DB().SetMaxIdleConns(databaseSetting.MaxIdleConns)
	db.DB().SetMaxOpenConns(databaseSetting.MaxOpenConns)
	otgorm.AddGormCallbacks(db)
	return db, nil
}

func beforeQuery(scope *gorm.Scope) {
	scope.Set(SlowQueryStartTimeTag, time.Now()) // 记录查询开始时间
}

func afterQuery(scope *gorm.Scope) {
	endTime := time.Now()
	startTimeObj, ok := scope.Get(SlowQueryStartTimeTag)
	if !ok {
		global.Logger.Error(context.Background(), "after query get start time err")
		return
	}

	startTime, ok := startTimeObj.(time.Time)
	if !ok {
		global.Logger.Error(context.Background(), "after query build startTimeObj to time.Time error")
		return
	}

	duration := endTime.Sub(startTime)
	if duration > global.DatabaseSetting.SlowThresholdDuration {
		slowSQLStr := fmt.Sprintf("[SLOW_SQL] %s duration: %v, sql: [%s] \n",
			currentFunction(4), duration, scope.SQL)
		if len(slowSQLStr) > 3000 {
			slowSQLStr = slowSQLStr[0:3000] // 超长时截断
		}
		global.Logger.Error(context.Background(), slowSQLStr)
		return
	}
}

// 获取调用堆栈信息
func currentFunction(skip int) string {
	pc, _, line, _ := runtime.Caller(skip)
	funcName := runtime.FuncForPC(pc)
	return fmt.Sprintf(" [method: %s, line: %d] ", funcName, line)
}

func updateTimeStampForCreateCallback(scope *gorm.Scope) {
	if !scope.HasError() {
		nowTime := time.Now().Unix()
		if createTimeField, ok := scope.FieldByName("CreatedOn"); ok {
			if createTimeField.IsBlank {
				_ = createTimeField.Set(nowTime)
			}
		}

		if modifyTimeField, ok := scope.FieldByName("ModifiedOn"); ok {
			if modifyTimeField.IsBlank {
				_ = modifyTimeField.Set(nowTime)
			}
		}
	}
}

func updateTimeStampForUpdateCallback(scope *gorm.Scope) {
	if _, ok := scope.Get("gorm:update_column"); !ok {
		_ = scope.SetColumn("ModifiedOn", time.Now().Unix())
	}
}

func deleteCallback(scope *gorm.Scope) {
	if !scope.HasError() {
		var extraOption string
		if str, ok := scope.Get("gorm:delete_option"); ok {
			extraOption = fmt.Sprint(str)
		}

		deletedOnField, hasDeletedOnField := scope.FieldByName("DeletedOn")
		isDelField, hasIsDelField := scope.FieldByName("IsDel")
		if !scope.Search.Unscoped && hasDeletedOnField && hasIsDelField {
			now := time.Now().Unix()
			scope.Raw(fmt.Sprintf(
				"UPDATE %v SET %v=%v,%v=%v%v%v",
				scope.QuotedTableName(),
				scope.Quote(deletedOnField.DBName),
				scope.AddToVars(now),
				scope.Quote(isDelField.DBName),
				scope.AddToVars(1),
				addExtraSpaceIfExist(scope.CombinedConditionSql()),
				addExtraSpaceIfExist(extraOption),
			)).Exec()
		} else {
			scope.Raw(fmt.Sprintf(
				"DELETE FROM %v%v%v",
				scope.QuotedTableName(),
				addExtraSpaceIfExist(scope.CombinedConditionSql()),
				addExtraSpaceIfExist(extraOption),
			)).Exec()
		}
	}
}

func addExtraSpaceIfExist(str string) string {
	if str != "" {
		return " " + str
	}
	return ""
}
