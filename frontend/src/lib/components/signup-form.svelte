<script lang="ts">
    import { goto } from '$app/navigation';
    import { Button } from "$lib/components/ui/button/index.js";
    import * as Card from "$lib/components/ui/card/index.js";
    import { Input } from "$lib/components/ui/input/index.js";
    import type { ComponentProps } from "svelte";
    import * as Field from "$lib/components/ui/field/index.js";
    import { register, loading, error, clearError } from '$lib/stores/auth';

    let { ...restProps }: ComponentProps<typeof Card.Root> = $props();

    // Form state
    let username = $state('');
    let password = $state('');
    let confirmPassword = $state('');
    let validationError = $state('');

    // Handle form submit
    async function handleSubmit(event: Event) {
        event.preventDefault();
        clearError();
        validationError = '';

        // Validasi password match
        if (password !== confirmPassword) {
            validationError = 'Password tidak cocok';
            return;
        }

        // Validasi password minimum 8 karakter
        if (password.length < 8) {
            validationError = 'Password minimal 8 karakter';
            return;
        }

        const success = await register(username, password);
        if (success) {
            goto('/wallet');
        }
    }
</script>

<Card.Root {...restProps}>
    <Card.Header>
        <Card.Title>Buat Akun</Card.Title>
        <Card.Description>Masukkan informasi untuk membuat akun baru</Card.Description>
    </Card.Header>
    <Card.Content>
        <form onsubmit={handleSubmit}>
            <Field.Group>
                {#if $error || validationError}
                    <div class="rounded-md bg-red-50 p-3 text-sm text-red-600">
                        {$error || validationError}
                    </div>
                {/if}
                <Field.Field>
                    <Field.Label for="name">Username</Field.Label>
                    <Input 
                        id="name" 
                        type="text" 
                        required 
                        bind:value={username}
                        disabled={$loading}
                    />
                </Field.Field>
                <Field.Field>
                    <Field.Label for="password">Password</Field.Label>
                    <Input 
                        id="password" 
                        type="password" 
                        required 
                        bind:value={password}
                        disabled={$loading}
                    />
                    <Field.Description>Minimal 8 karakter.</Field.Description>
                </Field.Field>
                <Field.Field>
                    <Field.Label for="confirm-password">Konfirmasi Password</Field.Label>
                    <Input 
                        id="confirm-password" 
                        type="password" 
                        required 
                        bind:value={confirmPassword}
                        disabled={$loading}
                    />
                    <Field.Description>Masukkan ulang password Anda.</Field.Description>
                </Field.Field>
                <Field.Group>
                    <Field.Field>
                        <Button type="submit" disabled={$loading}>
                            {$loading ? 'Memproses...' : 'Buat Akun'}
                        </Button>
                        <Field.Description class="px-6 text-center">
                            Sudah punya akun? <a href="/login">Masuk</a>
                        </Field.Description>
                    </Field.Field>
                </Field.Group>
            </Field.Group>
        </form>
    </Card.Content>
</Card.Root>
