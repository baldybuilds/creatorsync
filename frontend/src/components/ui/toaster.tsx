'use client';

import * as React from 'react';
import { Toast, ToastProvider, ToastViewport } from '../toaster';

export function Toaster() {
    return (
        <ToastProvider>
            <ToastViewport />
        </ToastProvider>
    );
}