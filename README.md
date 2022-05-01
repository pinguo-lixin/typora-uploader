# Typora Uploader

Typora 文件上传组件。用于将 typora 中的文件上传至七牛。

## 配置

1. 通过环境变量配置七牛 AK, SK

```bash
export QINIU_ACCESS_KEY=<access key>
export QINIU_SECRET_KEY=<secret key>
```

2. 配置 typora

路径：偏好设置 -> 图像

![image-20220501101718253](https://shine-doc.lixinshow.top/typora/2022/0a0883e37994f46145d24edbc7ced43f)

其中的命令：

`typora-uploader -bucket=<bucket> -host=<host> --`
