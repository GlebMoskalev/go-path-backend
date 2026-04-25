#!/bin/sh
set -e

if [ -f go.mod ] && grep -q "github.com/jackc/pgx" go.mod 2>/dev/null; then
    su-exec postgres pg_ctl -D /var/lib/postgresql/data -l /tmp/pg.log -w -t 10 start >/dev/null 2>&1 || {
        echo "FATAL: не удалось запустить PostgreSQL" >&2
        cat /tmp/pg.log >&2 2>/dev/null || true
        exit 1
    }

    for i in 1 2 3 4 5 6 7 8 9 10; do
        if su-exec postgres pg_isready -h localhost -p 5432 -q 2>/dev/null; then
            break
        fi
        sleep 0.3
    done
fi

exec "$@"
