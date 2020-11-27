package com.guigou.botisgud.commands

import com.gitlab.kordlib.common.entity.Snowflake
import com.gitlab.kordlib.core.Kord
import com.gitlab.kordlib.core.behavior.edit
import com.guigou.botisgud.extensions.logger
import com.guigou.botisgud.models.AbsoluteReminderTrigger
import com.guigou.botisgud.services.WordService
import com.guigou.botisgud.services.WordServiceImpl
import io.ktor.util.*
import kotlinx.coroutines.CoroutineScope
import kotlinx.coroutines.Dispatchers
import kotlinx.coroutines.delay
import kotlinx.coroutines.launch
import java.time.Instant.now

data class NicknameOptions(val commandTriggerExpression: String = "0 22 * * *")

class Nickname(
    private val users: Map<Snowflake, List<Snowflake>>,
    private val options: NicknameOptions = NicknameOptions(),
    private val wordService: WordService = WordServiceImpl(),
) : Command {
    companion object {
        val logger = logger()
    }

    override suspend fun register(client: Kord, scope: CoroutineScope) {
        logger.info("register nickname command")

        val trigger = AbsoluteReminderTrigger(options.commandTriggerExpression)
        scope.launch(Dispatchers.Default) {
            while (true) {
                val delayTill = trigger.timestamp()
                logger.info("delay processing till [${delayTill}]")
                delay(delayTill.toEpochMilli() - now().toEpochMilli())
                logger.info("begin processing")

                for (entry in users) {
                    for (userSnowflake in entry.value) {
                        try {
                            val member = client.getUser(userSnowflake)!!.asMember(entry.key)
                            logger.trace("processing [${member.displayName}]")
                            member.edit {
                                nickname = wordService.random()
                            }
                        } catch (ex: Exception) {
                            logger.error(ex)
                        }
                    }
                }

                logger.info("processing completed")
            }
        }
    }

}
