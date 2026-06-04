FROM oven/bun:1-slim AS base
WORKDIR /app

FROM base AS deps
COPY package.json bun.lock ./
RUN bun install --frozen-lockfile --production

FROM base
WORKDIR /app
COPY --from=deps /app/node_modules ./node_modules
COPY main.ts tsconfig.json readme.md ./
COPY src ./src

ENV NODE_ENV=production
ENV BUN_GARBAGE_COLLECTOR_LEVEL=2
ENV PORT=3000
EXPOSE 3000

CMD ["bun", "run", "--smol", "main.ts"]
