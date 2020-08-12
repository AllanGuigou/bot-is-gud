package com.guigou.botisgud.commands

import com.gitlab.kordlib.common.entity.Snowflake
import com.gitlab.kordlib.core.Kord
import com.gitlab.kordlib.core.entity.ReactionEmoji
import com.gitlab.kordlib.core.event.message.ReactionAddEvent
import com.gitlab.kordlib.core.supplier.EntitySupplier
import com.guigou.botisgud.services.ReminderService
import io.mockk.*
import kotlinx.coroutines.runBlocking
import org.junit.Test

class ReminderTests {

    private val serviceMock = mockk<ReminderService>()

    @Test
    fun `onReactionAddEvent adds a reminder when given a valid emoji`() = runBlocking {
        every { serviceMock.add(any(), any()) }
        val eventMock = mockk<ReactionAddEvent>()
        print("TEST")
        every { eventMock.emoji } answers {ReactionEmoji.Unicode("⌚")}
        every { eventMock.userId } answers {Snowflake("1")}
        coEvery { eventMock.message.asMessage().content } answers {"foo"}
        every { eventMock.guildId } answers  {Snowflake("2")}
        every { eventMock.channelId } answers {Snowflake("3")}
        every { eventMock.messageId } answers {Snowflake("4")}
        val sut = Reminder(serviceMock)

        sut.onReactionAddEvent(eventMock)

        print("TEST")
        verify { serviceMock.add(any(), any()) }
    }

    @Test
    fun `onReactionAddEvent does not add a reminder when given a invalid emoji`() = runBlocking {
        every { serviceMock.add(any(), any()) }
        val event = ReactionAddEvent(
                userId = Snowflake("1"),
                channelId = Snowflake("2"),
                messageId = Snowflake("3"),
                guildId = Snowflake("4"),
                emoji = ReactionEmoji.Unicode("⌚"),
                kord = mockk(),
                shard = 1,
                supplier = mockk()
        )
        val sut = Reminder(serviceMock)

        sut.onReactionAddEvent(event)

        verify { serviceMock.add(any(), any()) wasNot called }
    }
}
