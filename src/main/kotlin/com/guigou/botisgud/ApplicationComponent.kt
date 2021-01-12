package com.guigou.botisgud

import com.guigou.botisgud.commands.Nickname
import com.guigou.botisgud.commands.NicknameOptions
import com.guigou.botisgud.commands.Reminder
import com.guigou.botisgud.commands.Typing
import com.guigou.botisgud.services.nickname.NicknameService
import com.guigou.botisgud.services.nickname.NicknameServiceImpl
import com.guigou.botisgud.services.reminder.ReminderService
import com.guigou.botisgud.services.reminder.ReminderServiceImpl
import com.guigou.botisgud.services.user.UserService
import com.guigou.botisgud.services.user.UserServiceImpl
import com.guigou.botisgud.services.word.WordService
import com.guigou.botisgud.services.word.WordServiceImpl
import me.tatarka.inject.annotations.Component
import me.tatarka.inject.annotations.Provides

@Component
abstract class ApplicationComponent {
    abstract val nicknameCommand: Nickname
    abstract val reminderCommand: Reminder
    abstract val typingCommand: Typing

    abstract val reminderServiceImpl: ReminderServiceImpl

    @Provides
    protected fun reminderService(): ReminderService = reminderServiceImpl

    @Provides
    protected fun userService(): UserService = UserServiceImpl()

    @Provides
    protected fun nicknameOptions() = System.getenv("NICKNAME_COMMAND_TRIGGER_EXPRESSION").let {
        if (it.isNullOrEmpty()) NicknameOptions() else NicknameOptions(it)
    }

    @Provides
    protected fun nicknameService(): NicknameService = NicknameServiceImpl()

    @Provides
    protected fun wordService(): WordService = WordServiceImpl()
}
