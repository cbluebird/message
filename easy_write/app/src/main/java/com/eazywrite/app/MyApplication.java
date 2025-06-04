package com.eazywrite.app;

import android.app.Application;

import com.tencent.bugly.crashreport.CrashReport;

import org.litepal.LitePal;

public class MyApplication extends Application {

    public static MyApplication instance;

    public static MyApplication getInstance() {
        return instance;
    }

    @Override
    public void onCreate() {
        super.onCreate();
        instance = this;
        CrashReport.initCrashReport(this, "581128a204", BuildConfig.DEBUG);
        LitePal.initialize(this);
    }
}
