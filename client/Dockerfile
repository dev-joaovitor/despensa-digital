FROM node:24-alpine AS dev

WORKDIR /app

ENV COREPACK_ENABLE_DOWNLOADS=1

RUN corepack enable && corepack prepare pnpm@11.8.0 --activate

ENV PNPM_HOME="/pnpm"
ENV PATH="$PNPM_HOME/bin:$PATH"
ENV CI=true

COPY package*.json pnpm-lock.yaml pnpm-workspace.yaml* ./

# utilize BuildKit caching engine directly to prevent bloating docker image layers with store files
RUN --mount=type=cache,id=pnpm,target=/pnpm/store pnpm install --frozen-lockfile

COPY . .

EXPOSE 5173

CMD ["pnpm", "dev", "--host"]

FROM node:24-alpine AS production

RUN addgroup -S appgroup && adduser -S appuser -G appgroup

WORKDIR /app

ENV COREPACK_ENABLE_DOWNLOADS=1
RUN corepack enable && corepack prepare pnpm@11.8.0 --activate

ENV PNPM_HOME="/pnpm"
ENV PATH="$PNPM_HOME/bin:$PATH"
ENV CI=true

COPY package*.json pnpm-lock.yaml pnpm-workspace.yaml* ./

RUN --mount=type=cache,id=pnpm,target=/pnpm/store pnpm install --prod --frozen-lockfile

COPY --from=dev --chown=appuser:appgroup /app/build ./build

USER appuser

EXPOSE 3000

CMD ["pnpm", "build"]
