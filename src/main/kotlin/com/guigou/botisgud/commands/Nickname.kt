package com.guigou.botisgud.commands

import com.gitlab.kordlib.core.Kord
import com.gitlab.kordlib.core.behavior.edit
import com.guigou.botisgud.extensions.logger
import com.guigou.botisgud.models.AbsoluteReminderTrigger
import com.guigou.botisgud.services.nickname.NicknameService
import com.guigou.botisgud.services.nickname.NicknameServiceImpl
import com.guigou.botisgud.services.word.WordService
import com.guigou.botisgud.services.word.WordServiceImpl
import io.ktor.util.*
import kotlinx.coroutines.CoroutineScope
import kotlinx.coroutines.Dispatchers
import kotlinx.coroutines.delay
import kotlinx.coroutines.launch
import me.tatarka.inject.annotations.Inject
import java.time.Instant.now

data class NicknameOptions(val commandTriggerExpression: String = "0 22 * * *")

@Inject
class Nickname(
    private val options: NicknameOptions,
    private val nicknameService: NicknameService,
    private val wordService: WordService,
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

                val usersByGuild = nicknameService.getUsersByGuild()
                for (guild in usersByGuild) {
                    for (userSnowflake in guild.value) {
                        try {
                            val member = client.getUser(userSnowflake)!!.asMember(guild.key)
                            logger.trace("processing [${member.username}]")
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
