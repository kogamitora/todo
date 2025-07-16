# Load environment variables from .env file
if [ -f .env ]; then
    export $(cat .env | sed 's/#.*//g' | xargs)
fi

if ! nc -z localhost $PORT; then
    echo "Service is not running on port $PORT"
    exit 1
fi

echo "Service is running on port $PORT"
exit 0
