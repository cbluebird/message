package com.eazywrite.app.common

import android.content.Context
import android.widget.Toast
import android.widget.Toast.LENGTH_SHORT
import androidx.annotation.StringRes
import com.eazywrite.app.MyApplication
import java.lang.ref.WeakReference

// 使用WeakReference来持有Toast对象，避免内存泄漏
private var toast: WeakReference<Toast>? = null

/**
 * 显示一个Toast消息，使用资源ID来获取字符串。
 *
 * @param resId 要显示的字符串的资源ID。
 * @param context 用于获取字符串的上下文。
 * @param duration Toast显示的持续时间，默认为 LENGTH_SHORT。
 */
fun toast(
    @StringRes resId: Int,
    context: Context = MyApplication.instance,
    duration: Int = LENGTH_SHORT
) {
    Toast.makeText(context, resId, duration).show() // 使用资源ID创建Toast并显示
}

/**
 * 显示一个Toast消息，使用字符串。
 *
 * @param text 要显示的文本。
 * @param context 用于显示Toast的上下文。
 * @param duration Toast显示的持续时间，默认为 LENGTH_SHORT。
 */
@JvmOverloads
fun toast(text: String, context: Context = MyApplication.instance, duration: Int = LENGTH_SHORT) {
    // 如果之前有Toast正在显示，先取消它
    toast?.get()?.cancel()
    // 创建一个新的Toast对象，并使用WeakReference持有它，以避免内存泄漏
    toast = WeakReference(Toast.makeText(context.applicationContext, text, duration))
    // 显示Toast
    toast?.get()?.show()
}