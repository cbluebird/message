@file:OptIn(ExperimentalMaterial3Api::class, ExperimentalMaterial3Api::class)

package com.eazywrite.app.ui.importbill

import android.content.Context
import android.content.Intent
import android.net.Uri
import android.os.Bundle
import android.provider.Settings
import android.util.Log
import android.view.accessibility.AccessibilityNodeInfo
import android.widget.Toast
import androidx.activity.ComponentActivity
import androidx.activity.compose.setContent
import androidx.activity.viewModels
import androidx.compose.foundation.Image
import androidx.compose.foundation.layout.*
import androidx.compose.foundation.rememberScrollState
import androidx.compose.foundation.verticalScroll
import androidx.compose.material3.*
import androidx.compose.runtime.*
import androidx.compose.ui.Modifier
import androidx.compose.ui.layout.ContentScale
import androidx.compose.ui.platform.LocalContext
import androidx.compose.ui.res.painterResource
import androidx.compose.ui.unit.dp
import androidx.core.os.bundleOf
import com.eazywrite.app.R
import com.eazywrite.app.common.toast
import com.eazywrite.app.service.AutoAccessibilityService
import com.eazywrite.app.ui.theme.EazyWriteTheme
import com.eazywrite.app.util.*
import com.wilinz.accessbilityx.AccessibilityxService
import com.wilinz.accessbilityx.app.launchAppPackage
import com.wilinz.accessbilityx.bounds
import com.wilinz.accessbilityx.className1
import com.wilinz.accessbilityx.contentDescription1
import com.wilinz.accessbilityx.ensureClick
import com.wilinz.accessbilityx.text1
import kotlinx.coroutines.CoroutineScope
import kotlinx.coroutines.MainScope
import kotlinx.coroutines.delay
import kotlinx.coroutines.launch
/**
 * 支付宝账单导入Activity
 * 功能：提供支付宝账单文件选择和自动化操作支付宝获取账单
 */
@OptIn(ExperimentalMaterial3Api::class)
class AlipayImportActivity : ComponentActivity() {

    companion object {
        private const val TAG = "AlipayImportActivity"
    }
    // 使用viewModels委托初始化ViewModel
    private val viewModel: AlipayImportViewModel by viewModels()


    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)
        setWindow(isDarkStatusBarIcon = true) // 设置状态栏图标为深色
        setData(intent) // 处理传入的Intent数据
        setContent {
            EazyWriteTheme {
                Surface(
                    modifier = Modifier.fillMaxSize(),
                    color = MaterialTheme.colorScheme.background
                ) {
                    Content() // 主界面内容
                }
            }
        }
    }
    // 获取无障碍服务实例
    private val autoService: AutoAccessibilityService?
        get() = AutoAccessibilityService.instance
    // 协程作用域用于异步操作
    private val autoJob: CoroutineScope = MainScope()

    @Composable
    private fun Content() {
        Scaffold(topBar = {
            TopAppBar(
                title = { Text(text = "导入支付宝账单") }
            )
        }) {
            Column(
                modifier = Modifier
                    .padding(it)
                    .verticalScroll(rememberScrollState())
            ) {
                // 文件选择组件
                FileInput(onFileInput = {
                    viewModel.uri = it// 将选择的文件URI存入ViewModel
                })
                val context = LocalContext.current
                // 手动打开支付宝按钮
                ElevatedButton(
                    onClick = {
                        context.packageManager.getLaunchIntentForPackage("com.eg.android.AlipayGphone")
                            ?.let {
                                context.startActivity(it)
                            } ?: kotlin.run {
                            toast(text = "未安装支付宝")
                        }
                    },
                    modifier = Modifier
                        .padding(horizontal = 16.dp, vertical = 8.dp)
                        .fillMaxWidth()
                ) {
                    Text(text = "打开支付宝")
                }
                ElevatedButton(
                    onClick = {
                        autoJump(context)
                    },
                    modifier = Modifier
                        .padding(horizontal = 16.dp, vertical = 8.dp)
                        .fillMaxWidth()
                ) {
                    Text(text = "打开支付宝（自动操作）")
                }
                Import(viewModel)
                if (viewModel.uri != null) {
                    val filename = remember(viewModel.uri) {
                        viewModel.uri!!.getFileName(this@AlipayImportActivity)
                            ?: getRandomName()
                    }
                    Text("已选择：$filename", modifier = Modifier.padding(16.dp))
                } else {
                    Text("未选择任何文件", modifier = Modifier.padding(16.dp))
                }
                Column(
                    modifier = Modifier
                        .padding(horizontal = 16.dp)
                        .padding(bottom = 16.dp)
                ) {
                    Text("支付宝导入教程：", style = MaterialTheme.typography.titleLarge)

                }
                Image(
                    painter = painterResource(id = R.drawable.alipay_import),
                    modifier = Modifier.fillMaxWidth(),
                    contentDescription = "支付宝导入教程",
                    contentScale = ContentScale.FillWidth
                )
            }
        }
    }
    /**
     * 自动化操作支付宝流程
     * @param context 上下文对象
     */
    private fun autoJump(context: Context) {
        // 检查无障碍服务是否启用
        if (autoService == null || !AccessibilityxService.isAccessibilityServiceEnabled(
                context,
                AutoAccessibilityService::class.java
            )
        ) {
            toast(
                "请在无障碍设置中打开“${context.getString(R.string.text_import_invoices_automatically)}”",
                duration = Toast.LENGTH_LONG
            )
            // 跳转到无障碍设置界面
            context.startActivity(
                Intent(Settings.ACTION_ACCESSIBILITY_SETTINGS).addFlags(
                    Intent.FLAG_ACTIVITY_NEW_TASK
                )
            )
            return
        }
        // 尝试启动支付宝
        if (launchAppPackage("com.eg.android.AlipayGphone")) {
            autoJob.launch {
                autoService?.untilFindOne { it.contentDescription1 == "搜索框" }
                //通过控件的无障碍描述定位
                    ?.ensureClick()

                val editText =
                    autoService?.untilFindOne { it.className1 == "android.widget.EditText" }
                // 在搜索框输入"账单"
                val arguments =
                    bundleOf(AccessibilityNodeInfo.ACTION_ARGUMENT_SET_TEXT_CHARSEQUENCE to "账单")
                editText?.performAction(AccessibilityNodeInfo.ACTION_SET_TEXT, arguments)

                listOf(
                    "搜索",
                    "记录您的每一笔资金往来",
                    "收支分析",
                    "更多",
                    "开具交易流水证明",
                    "申请",
                    "完成"
                ).forEach { text ->
                    when (text) {

                        "收支分析" -> {
                            autoService?.untilFindOne { it.text1 == text }
                        }

                        "更多" -> {
                            autoService?.untilFindOne { it.contentDescription1 == text }
                                ?.ensureClick()
                        }

                        "申请", "完成" -> {
                            autoService?.untilFindOne { it.text1 == text }
                                ?.let { node ->
                                    delay(200)
                                    autoService?.click(node.bounds)
                                }
                        }

                        else -> {
                            autoService?.untilFindOne { it.text1 == text }
                                ?.ensureClick()
                        }
                    }

                }
    //      val successFlag =
    //                                    autoService?.untilFindOne { it.text1 == "邮件发送申请已提交" }
    //                                if (successFlag != null) {
    //                                    toast("申请已提交")
    //                                    autoService?.untilFindOne { it.text1 == "完成" }?.ensureClick()
                // 返回操作
                autoService?.untilFindOne { it.text1 == "账单" }
                for (i in 0 until 1) {
                    delay(500)
                    autoService?.back()
                }
    //                                }
            }
        } else {
            toast(text = "未安装支付宝")
        }
    }

    override fun onNewIntent(intent: Intent?) {
        super.onNewIntent(intent)
        Log.d(TAG, "onNewIntent: ")
        intent?.let {
            setData(it)
        }
    }// 处理新的Intent

    /**
     * 处理传入的Intent数据
     * @param intent 包含文件数据的Intent
     */
    private fun setData(intent: Intent) {
        Log.d(TAG, "setData:action: ${intent.action}")
        Log.d(TAG, "setData:data: ${intent}")
        when (intent.action) {
            // 处理分享发送的文件
            Intent.ACTION_SEND -> {
                intent.getParcelableExtraCompat<Uri>(Intent.EXTRA_STREAM)?.let {
                    viewModel.uri = it
                }
            }
            // 处理视图打开的文件
            Intent.ACTION_VIEW -> {
                intent.data?.let {
                    viewModel.uri = it
                }
            }
        }
    }


}





