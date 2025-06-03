

## 作用：OCR识别
1. 识别图片中的文字；
2. 返回文字内容和文本所在4个坐标点；
3. 返回文字识别相关的得分；


## 文档：
addr: https://cnocr.readthedocs.io/zh-cn/stable/

#### 在线 Demo:
地址：https://huggingface.co/spaces/breezedeus/CnOCR-Demo 
国内镜像：https://hf.qhduan.com/spaces/breezedeus/CnOCR-Demo




## 安装



## 使用
curl --location 'http://0.0.0.0:8501/ocr' \
--form 'image=@"/OCR/images/1.jpeg"'