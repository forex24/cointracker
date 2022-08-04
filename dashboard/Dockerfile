FROM node:18-buster-slim AS builder
WORKDIR /dashboard
COPY . .

RUN yarn install
RUN yarn build

FROM node:18-buster-slim AS dist
WORKDIR /dashboard
COPY --from=builder /dashboard/build /dashboard/build
RUN yarn global add serve
CMD [ "serve", "-s" , "build"]