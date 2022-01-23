package com.guigou.botisgud.services.nickname

import dev.kord.common.entity.Snowflake
import com.guigou.botisgud.extensions.logger
import org.jetbrains.exposed.sql.transactions.experimental.newSuspendedTransaction

class NicknameServiceImpl : NicknameService {
    companion object {
        val logger = logger()
    }

    override suspend fun getUsersByGuild(): Map<Snowflake, List<Snowflake>> {
        return newSuspendedTransaction {
            NicknameEntity.all().also {
                logger.trace("get [${it.count()}] nickname users")
            }.groupBy({ Snowflake(it.guildId) }, { Snowflake(it.userId) })
        }
    }
}
