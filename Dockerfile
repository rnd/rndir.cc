FROM nginx:alpine
EXPOSE 80
COPY assets /usr/share/nginx/html/assets
COPY out.html /usr/share/nginx/html/index.html
COPY resume/resume.pdf /usr/share/nginx/html/resume.pdf
