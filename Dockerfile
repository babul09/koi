
FROM nginx:alpine

# Remove default Nginx static files
RUN rm -rf /usr/share/nginx/html/*

COPY index.html /usr/share/nginx/html/
COPY wasm_exec.js /usr/share/nginx/html/
COPY game.wasm /usr/share/nginx/html/

# Expose port 80 for web traffic
EXPOSE 80

