mkdir -p $APP_LOG_PATH  /data   && ln -s $APP_LOG_PATH /data/logs && ls -la /data
/app/bin/app -conf /app/configs/@ENV/