package com.guigou.botisgud.services.nickname

import com.gitlab.kordlib.common.entity.Snowflake
import com.guigou.botisgud.extensions.logger
import me.tatarka.inject.annotations.Inject
import org.jetbrains.exposed.sql.transactions.experimental.newSuspendedTransaction

@Inject
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
