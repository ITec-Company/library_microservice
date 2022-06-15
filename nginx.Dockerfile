FROM nginx
RUN rm /etc/nginx/nginx.conf
WORKDIR /app
COPY nginx.conf /etc/nginx/nginx.conf
