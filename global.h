#ifndef GLOBAL_H
#define GLOBAL_H
#include <QWidget>
#include <functional>
#include "QStyle"
#include <QRegularExpression>
/******************************************************************************
 *
 * @file       global.h
 * @brief      刷新 qss
 *
 * @author     Liao
 * @date       2024/06/27
 * @history
 *****************************************************************************/
extern std::function<void(QWidget*)> repolish;
/**
 * @brief The ReqId enum 表示请求的id
 */
enum ReqId{
    ID_GET_VARIFY_CODE = 1001, //获取验证码
    ID_REG_USER = 1002, //注册用户
};

enum ErrorCodes{
    SUCCESS = 0,
    ERR_JSON = 1, //Json解析失败
    ERR_NETWORK = 2,
};

enum Modules{
    REGISTERMOD = 0,
};
#endif // GLOBAL_H
