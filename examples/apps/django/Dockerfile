FROM python@sha256:f79617b5d48f102ee95e3b7f0337d3495a17e1f64d66e1fefa5dce93b93b45d9
EXPOSE 8080
WORKDIR /app 
COPY requirements.txt /app
RUN pip3 install -r requirements.txt --no-cache-dir
COPY . /app 
ENTRYPOINT ["python3"] 
CMD ["manage.py", "runserver", "0.0.0.0:8080"]
