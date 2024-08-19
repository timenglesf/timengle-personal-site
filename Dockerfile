# Use a specific version of Debian to ensure consistency
FROM --platform=linux/amd64 debian:stable-slim

RUN apt-get update && apt-get install -y ca-certificates


# Copy the blog binary into the container's working directory
COPY ./timengledev_blog ./timengledev_blog

# Ensure the blog binary has execute permissions (just in case)
RUN chmod +x ./timengledev_blog

# Specify the entrypoint to run the blog binary
CMD ["./timengledev_blog"]
