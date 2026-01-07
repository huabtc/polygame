# Polygame Mobile App

Flutter 跨平台移动应用（iOS + Android）

## 技术栈

- **框架**: Flutter 3.x
- **状态管理**: Provider
- **HTTP 客户端**: Dio
- **本地存储**: SharedPreferences

## 功能特性

- ✅ 用户注册/登录
- ✅ 市场浏览
- ✅ 实时交易
- ✅ 持仓管理
- ✅ 订单历史
- ✅ 个人资料

## 快速开始

### 前置要求

- Flutter SDK 3.0+
- Android Studio / Xcode

### 安装依赖

```bash
flutter pub get
```

### 运行应用

```bash
# Android
flutter run

# iOS
flutter run -d ios

# 指定设备
flutter run -d <device_id>
```

### 构建应用

```bash
# Android APK
flutter build apk

# Android App Bundle
flutter build appbundle

# iOS
flutter build ios
```

## 项目结构

```
lib/
├── models/           # 数据模型
├── services/         # API 服务
├── providers/        # 状态管理
├── screens/          # 页面
├── widgets/          # 组件
├── utils/            # 工具类
└── main.dart         # 入口文件
```

## 配置

在 `lib/utils/constants.dart` 中配置 API 地址：

```dart
const String API_BASE_URL = 'http://your-api-url/api/v1';
```

## 支持平台

- ✅ Android 5.0+
- ✅ iOS 12.0+

## 许可证

MIT License
