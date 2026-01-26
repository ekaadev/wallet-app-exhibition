<script lang="ts">
    import { goto } from '$app/navigation';
    import { Button } from '$lib/components/ui/button/index.js';
    import { Input } from '$lib/components/ui/input/index.js';
    import { Field, FieldDescription, FieldGroup, FieldLabel } from '$lib/components/ui/field/index.js';
    import * as Card from '$lib/components/ui/card/index.js';
    import { login, loading, error, clearError } from '$lib/stores/auth';

    // Generate a unique ID for form fields
    let id = Math.random().toString(36).substring(2, 9);

    // Form state
    let username = $state('');
    let password = $state('');

    // Handle form submit
    async function handleSubmit(event: Event) {
        event.preventDefault();
        clearError();

        if (!username || !password) {
            return;
        }

        const success = await login(username, password);
        if (success) {
            goto('/wallet');
        }
    }
</script>

<Card.Root class="mx-auto w-full max-w-sm">
    <Card.Header>
        <Card.Title class="text-2xl">Login</Card.Title>
        <Card.Description>Masukkan username untuk login ke akun Anda</Card.Description>
    </Card.Header>
    <Card.Content>
        <form onsubmit={handleSubmit}>
            <FieldGroup>
                {#if $error}
                    <div class="rounded-md bg-red-50 p-3 text-sm text-red-600">
                        {$error}
                    </div>
                {/if}
                <Field>
                    <FieldLabel for="username-{id}">Username</FieldLabel>
                    <Input 
                        id="username-{id}" 
                        type="text" 
                        placeholder="username" 
                        required 
                        bind:value={username}
                        disabled={$loading}
                    />
                </Field>
                <Field>
                    <div class="flex items-center">
                        <FieldLabel for="password-{id}">Password</FieldLabel>
                    </div>
                    <Input 
                        id="password-{id}" 
                        type="password" 
                        placeholder="********" 
                        required 
                        bind:value={password}
                        disabled={$loading}
                    />
                </Field>
                <Field>
                    <Button type="submit" class="w-full" disabled={$loading}>
                        {$loading ? 'Memproses...' : 'Login'}
                    </Button>
                    <FieldDescription class="text-center">
                        Belum punya akun? <a href="/signup">Daftar</a>
                    </FieldDescription>
                </Field>
            </FieldGroup>
        </form>
    </Card.Content>
</Card.Root>
