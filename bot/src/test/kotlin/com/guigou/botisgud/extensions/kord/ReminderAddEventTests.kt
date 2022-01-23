package com.guigou.botisgud.extensions.kord

import assertk.assertThat
import assertk.assertions.isEqualTo
import dev.kord.common.entity.Snowflake
import dev.kord.core.Kord
import dev.kord.core.entity.ReactionEmoji
import dev.kord.core.event.message.ReactionAddEvent
import dev.kord.core.supplier.EntitySupplier
import io.mockk.mockk
import org.junit.jupiter.api.Test

class ReminderAddEventTests {

    private val kordMock = mockk<Kord>()
    private val entitySupplierMock = mockk<EntitySupplier>()

    @Test
    fun `link extension property returns guild message url`() {
        val sut = ReactionAddEvent(
            userId = Snowflake("1"),
            channelId = Snowflake("2"),
            messageId = Snowflake("3"),
            guildId = Snowflake("4"),
            emoji = ReactionEmoji.Unicode("\uD83C\uDFD3"),
            kord = kordMock,
            shard = 1,
            supplier = entitySupplierMock
        )

        val result = sut.link

        assertThat(result.toString()).isEqualTo("https://discordapp.com/channels/4/2/3")
    }

    @Test
    fun `link extension property returns direct message url`() {
        val sut = ReactionAddEvent(
            userId = Snowflake("1"),
            channelId = Snowflake("2"),
            messageId = Snowflake("3"),
            guildId = null,
            emoji = ReactionEmoji.Unicode("\uD83C\uDFD3"),
            kord = kordMock,
            shard = 1,
            supplier = entitySupplierMock
        )

        val result = sut.link

        assertThat(result.toString()).isEqualTo("https://discordapp.com/channels/@me/2/3")
    }
}