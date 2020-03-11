version: 2
jobs:
  build_and_test:
    docker:
      - image: circleci/golang:1.13
    working_directory: /go/src/{{gitRepoPath}}
    steps:
      - checkout
      - setup_remote_docker:
          docker_layer_caching: true
      - add_ssh_keys
{{#envFileName}}
      - run:
          name: Add environment variables to a file
          command: cp {{#envFileSampleName}} {{envFileName}}
{{/envFileName}}
      - run:
          name: Start Containers
          command: docker-compose -f docker-compose.yml up -d
      - run:
          name: Wait for Server
          command: |
            chmod +x .circleci/wait-for-server-start.sh
            .circleci/wait-for-server-start.sh
      - run:
          name: Wait extra 10s to ensure database is seeded
          command: sleep 10
      - run:
          name: Run tests
          command: docker exec -it {{containerName}} go test ./...

workflows:
  version: 2
  build:
    jobs:
      - build_and_test