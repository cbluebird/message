@file:OptIn(  // 启用实验性API
    ExperimentalMaterial3Api::class,  // Material3实验性API
    ExperimentalMaterialApi::class,   // Material实验性API
    ExperimentalFoundationApi::class  // Foundation实验性API
)

package com.eazywrite.app.ui.chart

import android.annotation.SuppressLint
import android.os.Bundle
import androidx.activity.ComponentActivity
import androidx.activity.compose.setContent
import androidx.compose.foundation.ExperimentalFoundationApi
import androidx.compose.foundation.layout.*
import androidx.compose.foundation.rememberScrollState
import androidx.compose.foundation.text.selection.SelectionContainer
import androidx.compose.foundation.verticalScroll
import androidx.compose.material.ExperimentalMaterialApi
import androidx.compose.material.icons.Icons
import androidx.compose.material.icons.filled.ExpandMore
import androidx.compose.material.pullrefresh.PullRefreshIndicator
import androidx.compose.material.pullrefresh.PullRefreshState
import androidx.compose.material.pullrefresh.pullRefresh
import androidx.compose.material.pullrefresh.rememberPullRefreshState
import androidx.compose.material3.*
import androidx.compose.runtime.*
import androidx.compose.runtime.saveable.rememberSaveable
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.graphics.Color
import androidx.compose.ui.platform.LocalDensity
import androidx.compose.ui.res.stringResource
import androidx.compose.ui.unit.dp
import androidx.compose.ui.viewinterop.AndroidView
import androidx.lifecycle.viewmodel.compose.viewModel
import com.eazywrite.app.R
import com.eazywrite.app.common.toast
import com.eazywrite.app.ui.theme.EazyWriteTheme
import com.eazywrite.app.ui.wiget.MarkdowmTextView
import com.eazywrite.app.util.fillMaxSize
import com.eazywrite.app.util.setWindow
import com.eazywrite.app.util.toArgbHex
import com.eazywrite.app.util.toRgbHex
import com.github.aachartmodel.aainfographics.aachartcreator.AAChartModel
import com.github.aachartmodel.aainfographics.aachartcreator.AAChartType
import com.github.aachartmodel.aainfographics.aachartcreator.AAChartView
import com.github.aachartmodel.aainfographics.aaoptionsmodel.AAStyle
import kotlinx.coroutines.Job
import kotlinx.coroutines.launch
import java.time.LocalDate

/**
 * 图表展示的Activity类，负责显示账单数据的图表和分析报告
 */
class ChartActivity : ComponentActivity() {
    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)
        // 设置窗口样式
        setWindow()
        // 设置Compose内容
        setContent {
            // 应用自定义主题
            EazyWriteTheme {
                // 创建背景容器
                Surface(color = MaterialTheme.colorScheme.background) {
                    // 显示图表页面
                    ChartPage()
                }
            }
        }
    }
}

// 日志标签常量
const val TAG = "ChartActivity.kt"

/**
 * 图表页面的主Composable函数
 */
@SuppressLint("UnusedMaterial3ScaffoldPaddingParameter")
@Composable
fun ChartPage() {
    // 获取ViewModel实例
    val vm: ChartViewModel = viewModel()

    // 当前选中的标签页索引（0=年视图，1=月视图）
    var tab by rememberSaveable { mutableStateOf(0) }

    // 控制年份选择器的显示状态
    var showYearPicker by rememberSaveable { mutableStateOf(false) }

    // 控制年月选择器的显示状态
    var showYearMonthPicker by rememberSaveable { mutableStateOf(false) }

    // 年视图选中的年份
    var yearOfYearView by rememberSaveable { mutableStateOf(LocalDate.now().year) }

    // 月视图选中的年份
    var yearOfMonthView by rememberSaveable { mutableStateOf(LocalDate.now().year) }

    // 月视图选中的月份
    var monthOfMonthView by rememberSaveable { mutableStateOf(LocalDate.now().month.value) }

    // 年份选择器弹窗
    if (showYearPicker) {
        com.eazywrite.app.ui.wiget.YearPicker(
            currentYear = yearOfYearView,
            onSelected = { selectedYear ->
                showYearPicker = false
                yearOfYearView = selectedYear
                vm.refresh(yearOfYearView) // 刷新年视图数据
            },
            onDismissRequest = { showYearPicker = false }
        )
    }

    // 年月选择器弹窗
    if (showYearMonthPicker) {
        com.eazywrite.app.ui.wiget.YearMonthPicker(
            currentYear = yearOfMonthView,
            currentMonth = monthOfMonthView,
            onSelected = { year, month ->
                showYearMonthPicker = false
                yearOfMonthView = year
                monthOfMonthView = month
                vm.refresh(yearOfMonthView, monthOfMonthView) // 刷新月视图数据
            },
            onDismissRequest = { showYearMonthPicker = false }
        )
    }

    /**
     * 刷新图表数据
     */
    fun refresh() {
        if (tab == 0) {
            vm.refresh(yearOfYearView) // 刷新年视图
        } else {
            vm.refresh(yearOfMonthView, monthOfMonthView) // 刷新月视图
        }
    }

    // 监听账单ID变化，自动刷新数据
    LaunchedEffect(key1 = Unit) {
        vm.getAllId().collect {
            refresh()
        }
    }

    // 使用Scaffold构建页面布局
    Scaffold(
        topBar = {
            // 顶部应用栏
            TopAppBar(
                title = {
                    Text(text = stringResource(id = R.string.chart)) // 标题
                },
                actions = {
                    // 报告相关状态
                    var showReportDialog by rememberSaveable { mutableStateOf(false) }
                    var isLoading by rememberSaveable { mutableStateOf(false) }
                    var job = remember<Job?> { null } // 报告生成任务

                    /**
                     * 生成报告
                     */
                    fun generateReports() {
                        if (vm.report != null) {
                            // 如果已有报告，直接显示
                            showReportDialog = true
                        } else {
                            // 开始生成新报告
                            toast("正在生成报告")
                            val year = if (tab == 0) yearOfYearView else yearOfMonthView
                            val month = if (tab == 0) null else monthOfMonthView
                            isLoading = true
                            // 启动报告生成任务
                            job = vm.generateReports1(
                                year,
                                month,
                                onSuccess = {
                                    showReportDialog = true // 成功时显示对话框
                                },
                                onFailure = {
                                    toast("失败：${it.message}") // 失败时提示
                                },
                                onComplete = {
                                    isLoading = false // 完成后重置状态
                                }
                            )
                        }
                    }

                    /**
                     * 处理报告生成请求
                     */
                    fun generateReports2() {
                        if (!isLoading) {
                            generateReports() // 开始生成报告
                        } else {
                            toast("已取消生成报告")
                            job?.cancel() // 取消正在进行的任务
                        }
                    }

                    // 报告显示对话框
                    if (showReportDialog && vm.report != null) {
                        AlertDialog(
                            onDismissRequest = { showReportDialog = false },
                            title = {},
                            text = {
                                // 可滚动的报告内容区域
                                Column(Modifier.verticalScroll(rememberScrollState())) {
                                    // 支持文本选择
                                    SelectionContainer {
                                        // 渲染Markdown格式的报告
                                        MarkdowmTextView(text = vm.report ?: "")
                                    }
                                }
                            },
                            confirmButton = {
                                // 确认按钮
                                TextButton(onClick = { showReportDialog = false }) {
                                    Text(text = "确定")
                                }
                            },
                            dismissButton = {
                                // 重新生成按钮
                                TextButton(onClick = {
                                    vm.report = null
                                    showReportDialog = false
                                    generateReports2()
                                }) {
                                    Text(text = "重新生成")
                                }
                            }
                        )
                    }

                    // 生成报告按钮
                    TextButton(onClick = { generateReports2() }) {
                        Text(
                            text = if (!isLoading) "生成报告" else "生成中（点击取消）",
                            style = MaterialTheme.typography.titleMedium
                        )
                    }

                    // 时间选择按钮（年视图）
                    if (tab == 0) {
                        TextButton(onClick = { showYearPicker = true }) {
                            Row(verticalAlignment = Alignment.CenterVertically) {
                                Text(
                                    text = "${yearOfYearView}年",
                                    style = MaterialTheme.typography.titleMedium
                                )
                                Icon(imageVector = Icons.Default.ExpandMore, contentDescription = null)
                            }
                        }
                    }
                    // 时间选择按钮（月视图）
                    else {
                        TextButton(onClick = { showYearMonthPicker = true }) {
                            Row(verticalAlignment = Alignment.CenterVertically) {
                                Text(
                                    text = "${yearOfMonthView}年${monthOfMonthView}月",
                                    style = MaterialTheme.typography.titleMedium
                                )
                                Icon(imageVector = Icons.Default.ExpandMore, contentDescription = null)
                            }
                        }
                    }
                }
            )
        }
    ) { pad ->  // Scaffold提供的内边距
        // 下拉刷新状态
        val refreshing by rememberSaveable { mutableStateOf(false) }
        val scope = rememberCoroutineScope()
        val pullRefreshState = rememberPullRefreshState(
            refreshing = refreshing,
            onRefresh = {
                scope.launch {
                    refresh() // 刷新数据
                }
            }
        )

        // 初始化时刷新数据
        LaunchedEffect(key1 = Unit) {
            refresh()
        }

        // 页面主要内容
        Column(
            Modifier
                .padding(pad)  // 应用Scaffold的内边距
                .fillMaxSize()
        ) {
            // 标签选择栏
            TabRow(selectedTabIndex = tab) {
                val tabGroup = remember { listOf("年视图", "月视图") }
                tabGroup.forEachIndexed { i, item ->
                    Tab(
                        selected = tab == i,
                        onClick = {
                            if (tab != i) {
                                tab = i
                                refresh() // 切换标签时刷新数据
                            }
                        },
                        text = { Text(text = item) },
                    )
                }
            }

            // 图表内容区域
            ChartContent(pullRefreshState, refreshing)
        }
    }
}

/**
 * 图表内容区域
 * @param pullRefreshState 下拉刷新状态
 * @param refreshing 是否正在刷新
 */
@Composable
private fun ChartContent(
    pullRefreshState: PullRefreshState,
    refreshing: Boolean
) {
    // 下拉刷新容器
    Box(modifier = Modifier.pullRefresh(pullRefreshState)) {
        // 可滚动的内容区域
        Column(
            modifier = Modifier
                .verticalScroll(rememberScrollState()) // 记住滚动位置
                .padding(16.dp) // 内边距
        ) {
            // 折线图容器
            Box(
                modifier = Modifier
                    .fillMaxWidth()
                    .height(350.dp)
            ) {
                LineChart1() // 显示折线图
            }

            Spacer(modifier = Modifier.padding(32.dp)) // 间距

            // 支出饼图容器
            Box(
                modifier = Modifier
                    .fillMaxWidth()
                    .height(250.dp)
            ) {
                PieChartOut() // 显示支出饼图
            }

            Spacer(modifier = Modifier.padding(32.dp)) // 间距

            // 收入饼图容器
            Box(
                modifier = Modifier
                    .fillMaxWidth()
                    .height(250.dp)
            ) {
                PieChartIn() // 显示收入饼图
            }

            Spacer(modifier = Modifier.height(100.dp)) // 底部间距
        }

        // 下拉刷新指示器
        PullRefreshIndicator(
            refreshing,
            pullRefreshState,
            Modifier.align(Alignment.TopCenter)
        )
    }
}

/**
 * 折线图组件
 */
@Composable
private fun LineChart1() {
    // 获取主题颜色
    val primary = MaterialTheme.colorScheme.primary
    val onBackground = MaterialTheme.colorScheme.onBackground
    val background = MaterialTheme.colorScheme.background

    // 获取当前密度下的像素值
    val width = with(LocalDensity.current) { 0.5.dp.toPx() }

    // 获取ViewModel
    val vm: ChartViewModel = viewModel()

    /**
     * 创建折线图模型
     */
    fun aaChartModel() = AAChartModel(
        legendEnabled = true, // 启用图例
        yAxisTitle = "", // Y轴标题
        dataLabelsEnabled = true, // 显示数据标签
        chartType = AAChartType.Line, // 折线图类型
        title = "趋势分析", // 图表标题
        backgroundColor = Color.Transparent.toArgbHex(), // 透明背景
        series = vm.lineChartData.toTypedArray(), // 图表数据系列
        categories = vm.lineChartXData.toTypedArray(), // X轴分类
        titleStyle = AAStyle().color(onBackground.toRgbHex()) // 标题样式
    )

    // 使用AndroidView嵌入原生图表组件
    AndroidView(
        factory = { context ->
            // 创建图表视图并填充布局
            AAChartView(context).fillMaxSize().apply {
                aa_drawChartWithChartModel(aaChartModel()) // 绘制图表
            }
        },
        update = { view ->
            // 更新图表
            view.aa_refreshChartWithChartModel(aaChartModel())
        }
    )
}

/**
 * 支出饼图组件
 */
@Composable
private fun PieChartOut() {
    // 获取主题颜色
    val primary = MaterialTheme.colorScheme.primary
    val onBackground = MaterialTheme.colorScheme.onBackground
    val background = MaterialTheme.colorScheme.background

    // 获取当前密度下的像素值
    val width = with(LocalDensity.current) { 0.5.dp.toPx() }

    // 获取ViewModel
    val vm: ChartViewModel = viewModel()

    /**
     * 创建饼图模型
     */
    fun aaChartModel() = AAChartModel(
        legendEnabled = false, // 禁用图例
        yAxisTitle = "消费分析", // Y轴标题
        dataLabelsEnabled = true, // 显示数据标签
        chartType = AAChartType.Pie, // 饼图类型
        title = "消费分析", // 图表标题
        backgroundColor = Color.Transparent.toArgbHex(), // 透明背景
        series = arrayOf(vm.pieChartOutData).mapNotNull { it }.toTypedArray(), // 图表数据
        titleStyle = AAStyle().color(onBackground.toRgbHex()) // 标题样式
    )

    // 使用AndroidView嵌入原生图表组件
    AndroidView(
        factory = { context ->
            AAChartView(context).fillMaxSize().apply {
                aa_drawChartWithChartModel(aaChartModel()) // 绘制图表
            }
        },
        update = { view ->
            view.aa_refreshChartWithChartModel(aaChartModel()) // 更新图表
        }
    )
}

/**
 * 收入饼图组件
 */
@Composable
private fun PieChartIn() {
    // 获取主题颜色
    val primary = MaterialTheme.colorScheme.primary
    val onBackground = MaterialTheme.colorScheme.onBackground
    val background = MaterialTheme.colorScheme.background

    // 获取当前密度下的像素值
    val width = with(LocalDensity.current) { 0.5.dp.toPx() }

    // 获取ViewModel
    val vm: ChartViewModel = viewModel()

    /**
     * 创建饼图模型
     */
    fun aaChartModel() = AAChartModel(
        legendEnabled = false, // 禁用图例
        yAxisTitle = "收入分析", // Y轴标题
        dataLabelsEnabled = true, // 显示数据标签
        chartType = AAChartType.Pie, // 饼图类型
        title = "收入分析", // 图表标题
        backgroundColor = Color.Transparent.toArgbHex(), // 透明背景
        series = arrayOf(vm.pieChartInData).mapNotNull { it }.toTypedArray(), // 图表数据
        titleStyle = AAStyle().color(onBackground.toRgbHex()) // 标题样式
    )

    // 使用AndroidView嵌入原生图表组件
    AndroidView(
        factory = { context ->
            AAChartView(context).fillMaxSize().apply {
                aa_drawChartWithChartModel(aaChartModel()) // 绘制图表
            }
        },
        update = { view ->
            view.aa_refreshChartWithChartModel(aaChartModel()) // 更新图表
        }
    )
}