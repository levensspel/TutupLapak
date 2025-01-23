package functionCallerInfo

type FunctionCaller string

const (
	UserControllerRegisterByEmail   FunctionCaller = "userController.RegisterByEmail"
	UserServiceRegisterByEmail      FunctionCaller = "userService.RegisterByEmail"
	UserRepositoryCreateUserByEmail FunctionCaller = "userRepository.CreateUserByEmail"

	UserControllerRegisterByPhone   FunctionCaller = "userController.RegisterByPhone"
	UserServiceRegisterByPhone      FunctionCaller = "userService.RegisterByPhone"
	UserRepositoryCreateUserByPhone FunctionCaller = "userRepository.CreateUserByPhone"

	UserControllerLoginByEmail   FunctionCaller = "userController.LoginByEmail"
	UserServiceLoginByEmail      FunctionCaller = "userService.LoginByEmail"
	UserRepositoryGetAuthByEmail FunctionCaller = "userRepository.GetAuthByEmail"

	UserControllerLoginByPhone   FunctionCaller = "userController.LoginByPhone"
	UserServiceLoginByPhone      FunctionCaller = "userService.LoginByPhone"
	UserRepositoryGetAuthByPhone FunctionCaller = "userRepository.GetAuthByPhone"
)
