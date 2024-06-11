FROM node:21 as base

ARG PORT=3000

ENV NODE_ENV=production

WORKDIR /src

# Build
FROM base as build

#TODO: add lock.json COPY --link package.json package-lock.json ./
COPY --link package.json ./
RUN npm install --production=false

COPY --link . .

# Any better ideas?
ENV NUXT_MAPBOX_ACCESS_TOKEN=pk.eyJ1Ijoia3VkMDAiLCJhIjoiY2x2N2g1eHF3MDVjcjJrcDZ0azhhZzhndCJ9.8kF5iT7DrylxHrArgEnXgQ
ENV NUXT_GITHUB_CLIENT_ID=d28501f357fa49a54cd5 
ENV NUXT_GITHUB_CLIENT_SECRET=fcecbd1b9ddee8e61d831d58cbe2d146a2ad1963
ENV AUTH_ORIGIN=http://localhost:3000

RUN npm run build
RUN npm prune

# Run
FROM base

ENV PORT=$PORT

COPY --from=build /src/.output /src/.output
# if rely on unbundled dependencies
# COPY --from=build /src/node_modules /src/node_modules

CMD [ "node", ".output/server/index.mjs" ]