pipeline {
    agent any
    environment {
        PORT =  credentials('BACKEND_PORT')

        DB_HOST = credentials('DB_HOST')
        DB_PORT = credentials('DB_PORT')
        DB_DATABASE = credentials('DB_DATABASE')
        DB_USERNAME = credentials('DB_USERNAME')
        DB_PASSWORD = credentials('DB_PASSWORD') 
        DB_ROOT_PASSWORD = credentials('DB_ROOT_PASSWORD')

        MAIL_HOST = credentials('MAIL_HOST')
        MAIL_PORT = credentials('MAIL_PORT')
        MAIL_USERNAME = credentials('MAIL_USERNAME')
        MAIL_PASSWORD = credentials('MAIL_PASSWORD')
        MAIL_FROM = credentials('MAIL_FROM')
        MAIL_FROM_NAME = credentials('MAIL_FROM_NAME')

        REDIS_HOST = credentials('REDIS_HOST')
        REDIS_PORT = credentials('REDIS_PORT')
        REDIS_PASSWORD = credentials('REDIS_PASSWORD')

        JWT_SECRET = credentials('JWT_SECRET') 
        
        ALLOWED_ORIGINS = credentials('ALLOWED_ORIGINS')
    }

    stages {
        stage('Remove Old Docker Image') {
            steps {
                script {
                    // Stop and remove the old container if it exists
                    sh 'docker stop travel-be-container || true'
                    sh 'docker rm travel-be-container || true'
                    
                    // Remove the old image
                    sh 'docker rmi $(docker images -q travel-be) || true'  // || true to prevent error if no image exists
                }
            }
        }

        stage('Build Docker Image') {
            steps {
                script {
                    // Build the new Docker image
                    sh 'docker build -t travel-be .'
                }
            }
        }

        stage('Run Container') {
            steps {
                script {
                    // Run the newly built container with environment variables
                    sh '''
                        docker run -d \
                            --restart unless-stopped \
                            --name travel-be-container \
                            -p ${PORT}:${PORT} \
                            -e PORT="$PORT" \
                            -e DB_HOST="$DB_HOST" \
                            -e DB_PORT="$DB_PORT" \
                            -e DB_DATABASE="$DB_DATABASE" \
                            -e DB_USERNAME="$DB_USERNAME" \
                            -e DB_PASSWORD="$DB_PASSWORD" \
                            -e DB_ROOT_PASSWORD="$DB_ROOT_PASSWORD" \
                            -e MAIL_HOST="$MAIL_HOST" \
                            -e MAIL_PORT="$MAIL_PORT" \
                            -e MAIL_USERNAME="$MAIL_USERNAME" \
                            -e MAIL_PASSWORD="$MAIL_PASSWORD" \
                            -e MAIL_FROM="$MAIL_FROM" \
                            -e MAIL_FROM_NAME="$MAIL_FROM_NAME" \
                            -e REDIS_HOST="$REDIS_HOST" \
                            -e REDIS_PORT="$REDIS_PORT" \
                            -e REDIS_PASSWORD="$REDIS_PASSWORD" \
                            -e JWT_SECRET="$JWT_SECRET" \
                            -e ALLOWED_ORIGINS="$ALLOWED_ORIGINS" \
                            travel-be
                    '''
                }
            }
        }
    }
}
