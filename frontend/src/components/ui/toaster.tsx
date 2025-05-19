'use client';

import * as React from 'react';
import { ToastProvider, ToastViewport } from '../toaster';

export function Toaster() {
    return (
        <ToastProvider>
            <ToastViewport />
        </ToastProvider>
    );
}